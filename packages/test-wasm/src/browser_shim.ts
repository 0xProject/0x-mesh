import { getSchemaValidator } from '@0x/mesh-browser-lite/lib/schema_validator';

(window as any).setSchemaValidator = (chainId: number, exchangeAddress: string, customOrderFilter: string) => {
    (window as any).schemaValidator = getSchemaValidator(chainId, exchangeAddress, JSON.parse(customOrderFilter));
};
