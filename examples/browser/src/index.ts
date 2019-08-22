import { Mesh, OrderEvent } from '@0x/mesh-browser';

(async () => {
    const mesh = new Mesh({
        ethereumRPCURL: 'https://mainnet.infura.io/v3/af2e590be00f463fbfd0b546784065ad',
        ethereumNetworkID: 1,
    });
    mesh.setOrderEventsHandler((events: Array<OrderEvent>) => {
        for (let event of events) {
            console.log(event);
        }
    });
    await mesh.startAsync();
})().catch(err => {
    console.error(err);
});
