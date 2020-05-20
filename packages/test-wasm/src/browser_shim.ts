import { getSchemaValidator } from '@0x/mesh-browser-lite/lib/schema_validator';
import { JsonSchema } from '@0x/mesh-browser-lite/lib/types';

(window as any).setSchemaValidator = (chainId: number, exchangeAddress: string, customOrderFilter: JsonSchema) => {
    (window as any).schemaValidator = getSchemaValidator(chainId, exchangeAddress, customOrderFilter);
};
