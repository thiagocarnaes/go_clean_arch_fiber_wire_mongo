// MongoDB initialization script
// This script runs when the MongoDB container starts for the first time

// Switch to the user_management database
db = db.getSiblingDB('user_management');

// Create collections with validation schemas
db.createCollection('users', {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: ["_id", "name", "email"],
         properties: {
            _id: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            name: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            email: {
               bsonType: "string",
               pattern: "^.+@.+\..+$",
               description: "must be a valid email address and is required"
            }
         }
      }
   }
});

db.createCollection('groups', {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: ["_id", "name", "members"],
         properties: {
            _id: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            name: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            members: {
               bsonType: "array",
               items: {
                  bsonType: "string"
               },
               description: "must be an array of strings and is required"
            }
         }
      }
   }
});

// Create indexes for better performance
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "name": 1 });
db.groups.createIndex({ "name": 1 });
db.groups.createIndex({ "members": 1 });

// Insert some sample data for development (optional)
db.users.insertMany([
   {
      _id: "user1",
      name: "John Doe",
      email: "john.doe@example.com"
   },
   {
      _id: "user2", 
      name: "Jane Smith",
      email: "jane.smith@example.com"
   },
   {
      _id: "user3",
      name: "Bob Johnson", 
      email: "bob.johnson@example.com"
   }
]);

db.groups.insertMany([
   {
      _id: "developers",
      name: "Development Team",
      members: ["user1", "user2"]
   },
   {
      _id: "managers",
      name: "Management Team", 
      members: ["user3"]
   },
   {
      _id: "all-staff",
      name: "All Staff",
      members: ["user1", "user2", "user3"]
   }
]);

print("Database initialization completed successfully!");
print("Created collections: users, groups"); 
print("Created indexes for performance optimization");
print("Inserted sample data for development");
