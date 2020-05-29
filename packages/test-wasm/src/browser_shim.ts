import { createSchemaValidator } from '@0x/mesh-browser-lite/lib/schema_validator';
import { createDatabase } from '@0x/mesh-browser-lite/lib/database';

(window as any).createSchemaValidator = createSchemaValidator;
(window as any).__mesh_dexie_newDatabase__ = createDatabase;
