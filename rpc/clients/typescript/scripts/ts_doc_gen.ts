import { logUtils, promisify } from '@0x/utils';
import * as fs from 'fs';
import * as glob from 'glob';
import { exec as execAsync } from 'promisify-child-process';
import * as rimraf from 'rimraf';
import * as yargs from 'yargs';

const rimrafAsync = promisify(rimraf);

(async () => {

    const args = yargs
    .option('sourceDir', {
        alias: ['s', 'src'],
        describe: 'Folder where the source TS files are located',
        type: 'string',
        normalize: true,
        demandOption: true,
    })
    .option('output', {
        alias: ['o', 'out'],
        describe: 'Folder where to put the output doc files',
        type: 'string',
        normalize: true,
        demandOption: true,
    })
    .example(
        "$0 --src 'src' --out 'docs'",
        'Full usage example',
    ).argv;

    await rimrafAsync(args.output);
    try {
        await execAsync(`./node_modules/typedoc/bin/typedoc --theme markdown --platform gitbook --excludePrivate --excludeProtected --excludeExternals --excludeNotExported --target ES5 --module commonjs --hideGenerator --out ${args.output} ${args.sourceDir}`);
    } catch (err) {
        logUtils.log('typedoc command failed: ', err);
        process.exit(1);
    }

    // Concat all TS Client MD files together into a single reference doc
    const referencePath = `${args.output}/reference.md`;
    await rimrafAsync(referencePath);
    glob(`${args.output}/**/*`, (err: Error | null, paths: string[]) => {
        if (err !== null) {
            throw err;
        }
        (paths as any).sort((firstEl: string, secondEl: string) => {
            const isFirstAFile = firstEl.includes('.md');
            const isSecondAFile = secondEl.includes('.md');
            if ((isFirstAFile && isSecondAFile) || (!isFirstAFile && !isSecondAFile)) {
                return 0;
            }
            if (isFirstAFile) {
                return -1;
            }
            if (isSecondAFile) {
                return 1;
            }
            return undefined;
        });
        for (const path of paths) {
            if (path.includes('.md', 1)) {
                if (!path.includes('README.md', 1) && !path.includes('/modules/', 1) && !path.includes('globals.md', 1)) {
                    // Read file content and concat to new file
                    const content = fs.readFileSync(path);
                    fs.appendFileSync(referencePath, content);
                    fs.appendFileSync(referencePath, '\n');
                    fs.appendFileSync(referencePath, '\n');
                    fs.appendFileSync(referencePath, '<hr />');
                    fs.appendFileSync(referencePath, '\n');
                    fs.appendFileSync(referencePath, '\n');
                }
                if (!path.includes('README.md', 1)) {
                    fs.unlinkSync(path);
                }
            } else {
                fs.rmdirSync(path);
            }
        }

        // Find/replace relative links with hash links
        const docsBuff = fs.readFileSync(referencePath);
        let docs = docsBuff.toString();
        docs = docs.replace(/\]\(((?!.*(github.com)).*)(#.*\))/g, ']($3');
        docs = docs.replace(/\]\(..\/interfaces\/.*?\.(.*?)\.md\)/g, '](#interface-$1)');
        docs = docs.replace(/\]\(..\/classes\/.*?\.(.*?)\.md\)/g, '](#class-$1)');
        docs = docs.replace(/\]\(..\/enums\/.*?\.(.*?)\.md\)/g, '](#enumeration-$1)');
        docs = docs.replace(/\]\(_types_\.(.*?)\.md\)/g, '](#interface-$1)');
        docs = docs.replace(/\]\(.*\.(.*?)\.md\)/g, '](#class-$1)');
        fs.writeFileSync(referencePath, docs);
        logUtils.log('TS doc generation complete!');
    });
})();
