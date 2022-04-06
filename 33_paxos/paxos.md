# paxos

## Proposal Numbers
一个 Proposal 是由 Proposal number 和 Proposal value 组成的，
首先来看 Proposal number，既然针对一个 Proposal，那么每个 Proposal 都必须有一个唯一的 ID，也就是 Proposal Number，
其有几个基本特性（注意后文不区分 Proposal number 和 Proposal Id）：

- 全局唯一，无论是相同的 Proposer 的不同 Proposal 还是，还是不同的 Proposer 的相同或者不同 Proposal
- 单调递增，越大优先级越高

## basic paxos
basic paxos 算法分为两个阶段，Prepare 和 Accept 阶段


### prepare阶段
1. Proposer发送propose
    1. Proposer 生成全局唯一且递增的Proposal ID（最简单的方法：round number + serverID），向集群的所有机器发送Propose，这里无需携带提案内容，只携带Proposal ID即可
2. Acceptor应答Propose
   1. Acceptor 收到Propose后，将会持久化proposal ID，并做出两个承诺，一个应答
      - 两个承诺
        1. 不再应答Proposal ID 小于等于（注意：这里是<= ）当前请求的Proposal
        2. 不再应答Proposal ID 小于（注意：这里是< ）当前请求的Accept请求
      - 一个应答
        1. 返回已经Accept 过的提案中Proposal ID 最大的那个提案的Value和accepted Proposal ID，没有则返回空值
        这应该也是为什么，在多个proposor的共识方案中，必须要有prepare阶段，因为如果没有，所有的acceptor都不知道哪个值是已经形成决议的。
3. proposer收到应答之后
   1. proposal至少要收到一半以上acceptor的response
      1. acceptor的response中只要存在一个acceptedValue，需要将自己提议的value换成收到的acceptedValue中，acceptedProposal中最大的那个
      无论response中的acceptedValue是否真的被系统chosen了。这点保证了无论如何，basic paxos系统是安全的。最终都能达成一致。

### accept阶段
1. Proposer 发送Accept
- 提案生成规则：Proposer 收集到多数派的Propose应答后，从应答中选择存在提案Value的并且同时也是Proposal ID最大的提案的Value，作为本次要发起Accept 的提案。如果所有应答的提案Value均为空值，则可以自己随意决定提案Value。然后携带上当前Proposal ID，向集群的所有机器发送Accept请求
2. 应答Accept
   - Acceptor 收到Accpet请求后，检查不违背自己之前作出的“两个承诺”情况下，持久化当前Proposal ID 和提案Value。最后Proposer 收集到多数派的Accept应答后，形成决议
     1. 持久化之后当前决议就已经形成，保证后续不会丢失
     2. 如果已经有决议决定了某个value，那么其他没有参与到这次决议的proposer会在后续的proposal中学习到这个value（也就是前文讲的回复Proposal ID最大的提案的Value，作为本次要发起Accept 的提案，这个学习过程在propose阶段就可以完成）
3. proposer处理应答
   - 如果acceptor没有接受自己的请求（response中的minProposal大于自己的proposalId），重新propose
     - 注意，basic paxos系统没有拒绝已经chosen的value再次proposal。只不过如果value已经chosen了，那么proposer在提议的过程中会学习到这个value，然后以更高的proposalId去提议。
   - 如果超过半数acceptor接受了，或者超过半数没接受但是minProposal一致，那么表明当前决议的value已经被paxos系统决定（chosen）了

### paxos算法的关键特性
#### safety
1. quorum 机制：因为两个半数必然有相交，这样就不可能到存在两个value被同时最终接受造成不一致。而且可用性更好，可以容忍半数以下机器失效
2. 后者认同前者：就是一旦一个acceptor接受一个确定性取值之后，后者会认同前者，不会进行破坏（这个取值其实是学习到的，认同前者，就是）。这样实现一种自动的往“某个提议值被多数派接受并生效”这一最终目标靠拢。
   - **这实际上也paxos核心解决的问题** 当一个提议被多数派接受后，这个提议对应的值被Chosen（选定），一旦有一个值被Chosen，那么只要按照协议的规则继续交互，后续被Chosen的值都是同一个值，也就是这个Chosen值的共识问题。
   - 一个acceptor只要接受了某个proposal，那么他就会认同这个值，并在proposal的response中返回。无论这个值是否真的被chosen了。

#### liveness
1. 不会死锁（抢占式）：在 Proposal 和 Accept 有两个阶段都有两个承诺，例如，大的 Proposal Id 会覆盖让前一个 Proposal accept 失败，这样有什么好处呢？这里其实就是为了保证不会出现死锁，如可以让大 Proposal Id 可以抢占前面小的 Proposal id propose 权利，继续进行，因为如果某acceptor在收到大多数的 Propose 回复之后挂了，不支持抢占的模式将会造成死锁
2. 会活锁：两个提案不断的出现 Prepare 阶段抢占。尽管出现的概率很小，但是从理论，确实存在活锁的可能性。最简单的解决方法是：proposer随机超时重试即可。当然，Multi-Paxos有leader，避免活锁。


## Multi-Paxos

Multi-Paxos中，所有的节点既是proposer，又是acceptor，外部client发给任意一个Multi-Paxos节点写日志的请求，由该节点负责propose该请求到Multi-Paxos集群中。
做出这种集群方式的假设只是为了方便后面的讲解（同时也是Ongaro视频使用的方法），不同的系统可以由不同的实现。

- Consensus Module：从Basic Paxos的决策单个Value，到连续决策多个Value。同时确保每个Servers上的顺序完全一致
- 相同顺序执行命令，所有Servers状态一致
- 只要多数节点存活，Server能提供正常服务