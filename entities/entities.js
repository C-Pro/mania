/**
 * Based on https://github.com/dialogflow/city-streets-trivia-nodejs
 * Copyright 2017 Google Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

'use strict';

// Loads in our node.js client library which we will use to make API calls
const dialogflow = require('dialogflow');

// Create a new EntityTypesClient, which
// communicates with the EntityTypes API endpoints
const entitiesClient = new dialogflow.EntityTypesClient();

// Create a path string for our agent based
// on its project ID (from first tab of Settings).
const projectId = 'mania-25a3b';
const agentPath = entitiesClient.projectAgentPath(projectId);

// Define an EntityType to represent categories.
const categoryEntityType = {
    displayName: 'category',
    kind: 'KIND_MAP',
    // List all of the Entities within this EntityType.
    entities: [
        { value: 'New York', synonyms: ['New York', 'NYC'] },
        { value: 'Los Angeles', synonyms: ['Los Angeles', 'LA', 'L.A.'] },
    ],
};

// Build a request object in the format the client library expects.
const categoryRequest = {
    parent: agentPath,
    entityType: categoryEntityType,
};

// Tell client library to call Dialogflow with
// a request to create an EntityType.
entitiesClient
    .createEntityType(categoryRequest)
    // Dialogflow will respond with details of the newly created EntityType.
    .then((responses) => {
        console.log('Created new entity type:', responses[0]["name"]);

        // Define and create an EntityType to represent menu items.
        const itemEntityType = {
            displayName: 'item',
            kind: 'KIND_MAP',

            entities: [
                { value: 'Broadway', synonyms: ['Broadway'] },
            ]
        };

        const itemRequest = {
            parent: agentPath,
            entityType: itemEntityType,
        };

        return entitiesClient.createEntityType(itemRequest);
    })
    // Dialogflow again responds with details of the newly created EntityType.
    .then((responses) => {
        console.log('Created new entity type:', responses[0]["name"]);
    })
    // Log any errors.
    .catch((err) => {
        console.error('Error creating entity type:', err);
    });
