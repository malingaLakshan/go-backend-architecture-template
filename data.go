LOC_MONGO_USER=$(sudo docker exec mongodb printenv MONGO_INITDB_ROOT_USERNAME)
LOC_MONGO_PASS=$(sudo docker exec mongodb printenv MONGO_INITDB_ROOT_PASSWORD)

sudo docker exec mongodb mongosh --quiet \
  --username "$LOC_MONGO_USER" \
  --password "$LOC_MONGO_PASS" \
  --authenticationDatabase admin \
  --eval '
const targetItem = "306AD5FF20600CD7AF69798A";
const targetFloor = "sim-5260-world";
const ignoredDatabases = new Set(["admin", "config", "local"]);

const databaseNames = db
  .getSiblingDB("admin")
  .adminCommand({listDatabases: 1})
  .databases
  .map(database => database.name)
  .filter(name => !ignoredDatabases.has(name));

let matchesFound = 0;
const collections = [];

print("SEARCHING ALL MONGODB DATABASES AND COLLECTIONS");
print("Target item_id: " + targetItem);
print("Target floor: " + targetFloor);

for (const databaseName of databaseNames) {
  const currentDatabase = db.getSiblingDB(databaseName);

  for (const collectionName of currentDatabase.getCollectionNames()) {
    const collection = currentDatabase.getCollection(collectionName);
    collections.push([databaseName, collectionName]);

    try {
      const query = {
        $or: [
          {item_id: targetItem},
          {itemId: targetItem},
          {tag_id: targetItem},
          {tagId: targetItem},
          {"location.item_id": targetItem},
          {"location.itemId": targetItem},
          {"locations.item_id": targetItem},
          {"locations.itemId": targetItem}
        ]
      };

      const documents = collection.find(query).limit(5).toArray();

      for (const document of documents) {
        print("\n==================================================");
        print("EXACT ITEM MATCH: " + databaseName + "." + collectionName);
        print("==================================================");
        printjson(document);
        matchesFound++;
      }
    } catch (error) {
    }
  }
}

if (matchesFound === 0) {
  print("\nExact item was not found. Searching recent documents by floor...");

  for (const [databaseName, collectionName] of collections) {
    if (matchesFound >= 5) {
      break;
    }

    try {
      const collection = db
        .getSiblingDB(databaseName)
        .getCollection(collectionName);

      const documents = collection
        .find({})
        .sort({_id: -1})
        .limit(1000)
        .toArray();

      for (const document of documents) {
        const text = EJSON.stringify(document);

        if (text.includes(targetFloor)) {
          print("\n==================================================");
          print("FLOOR MATCH: " + databaseName + "." + collectionName);
          print("==================================================");
          printjson(document);
          matchesFound++;

          if (matchesFound >= 5) {
            break;
          }
        }
      }
    } catch (error) {
    }
  }
}

if (matchesFound === 0) {
  print("\nNO LOCATION DOCUMENT FOUND.");
  print("Available collections and their sample field names:");

  for (const [databaseName, collectionName] of collections) {
    try {
      const document = db
        .getSiblingDB(databaseName)
        .getCollection(collectionName)
        .findOne({});

      if (document) {
        print(
          databaseName + "." + collectionName +
          " -> " + Object.keys(document).join(", ")
        );
      }
    } catch (error) {
    }
  }
} else {
  print("\nTOTAL MATCHING DOCUMENTS: " + matchesFound);
  print("Compare the named fields containing values 16 and 1.");
}
'

unset LOC_MONGO_USER LOC_MONGO_PASS