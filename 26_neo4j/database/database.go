package database

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
)

func CreateDriver(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
}

func CloseDriver(driver neo4j.Driver) error {
	return driver.Close()
}

func NodeCreate(driver neo4j.Driver, Cypher string, DB string) error {

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(Cypher, nil)
		if err != nil {
			log.Println("wirte to DB with error:", err)
			return nil, err
		}
		return result.Consume()
	})

	return err
}

/*
	NodeQuery is common API for Querying NODE in neo4j DB
   	the cypher string must use "n"(node) as the obj name
	example:
		  "MATCH (n) RETURN n, n.Desc as desString, n.name as objName LIMIT 100"
    neo4j.Node slice are return value.
*/
func NodeQuery(driver neo4j.Driver, Cypher string, DB string) ([]neo4j.Node, error) {

	var list []neo4j.Node
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(Cypher, nil)
		if err != nil {
			return nil, err
		}

		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("n"); ok {
				node := value.(neo4j.Node)
				list = append(list, node)
			}
		}
		if err = result.Err(); err != nil {
			return nil, err
		}

		return list, result.Err()
	})

	if err != nil {
		log.Println("Read error:", err)
	}
	return list, err
}

/*
	EdgeQuery is common API for Querying relationship in neo4j DB
   	the cypher string must use "r"(relationship) as the obj name
	example:
		  "MATCH (r) RETURN n, r.weight as weight LIMIT 100"
    neo4j.Relationship slice are return value.
*/
func EdgeQuery(driver neo4j.Driver, Cypher string, DB string) ([]neo4j.Relationship, error) {

	var list []neo4j.Relationship
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(Cypher, nil)
		if err != nil {
			log.Println("EdggeQuery Run failed: ", err)
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("r"); ok {
				relationship := value.(neo4j.Relationship)
				list = append(list, relationship)
				//				log.Println("Edgeid:", relationship.Id, ">>>Node:", relationship.StartId, "---", relationship.Type, "--->","Node:",relationship.EndId)
			}
		}
		if err = result.Err(); err != nil {
			return nil, err
		}
		return list, result.Err()
	})

	if err != nil {
		log.Println("Read error:", err)
	}
	return list, err
}
