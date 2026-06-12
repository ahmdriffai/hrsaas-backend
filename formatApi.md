- endpoint : /api/nasabah-kredit/search
  method : GET
  description : Search for items based on query parameters.
  query parameters :
   ```json
   {
     "query": "string", // The search term to look for in item names and descriptions.
     "category": "string", // Optional category to filter the search results.
     "limit": "integer" // Optional limit on the number of results returned (default is 10).
   }
   ```
  response :
   ```json
   {
     "results": [
       {
         "id": "string", // Unique identifier for the item.
         "name": "string", // Name of the item.
         "description": "string", // Description of the item.
         "category": "string" // Category of the item.
       }
     ],
     "totalResults": "integer" // Total number of items matching the search criteria.
   }
   ```