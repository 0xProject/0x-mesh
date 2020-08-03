import { MeshGraphQLClient } from './index';

(async () => {
    const client = new MeshGraphQLClient('http://localhost:60557/graphql', 'ws://localhost:60557/graphql');
    console.log('subscribing...');
    client.onOrderEvents().subscribe({
        start: () => {
            console.log('succesfully connected');
        },
        next: orderEvents => console.log(`received ${orderEvents.length} events`),
        error: err => {
            console.error(err);
            throw err;
        },
        complete: () => console.log('observable is complete'),
    });
    console.log('done subscribing');
    setInterval(() => console.log('tick'), 3000);
})().catch(err => console.error(err));
