import { Mesh } from '@0xorg/mesh-browser';

(async () => {
    const mesh = new Mesh({
        ethereumRPCURL: 'https://mainnet.infura.io/v3/af2e590be00f463fbfd0b546784065ad',
        ethereumNetworkID: 1,
    });
    await mesh.startAsync();
})().catch(err => {
    console.error(err);
});
