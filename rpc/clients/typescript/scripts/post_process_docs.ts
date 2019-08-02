import * as fs from 'fs';
import * as glob from 'glob';

function main(): void {
    // Concat all TS Client MD files together
    const tsClientDocsRoot = `${process.cwd()}/../../../docs/json_rpc_clients/typescript`;
    const allMdPath = `${tsClientDocsRoot}/all.md`;

    if (fs.existsSync(allMdPath)) {
        fs.unlinkSync(allMdPath);
    }
    glob(`${tsClientDocsRoot}/**/*`, (err: Error | null, paths: string[]) => {
        if (err !== null) {
            throw err;
        }
        (paths as any).sort((firstEl: string, secondEl: string) => {
            const firstIsFile = firstEl.includes('.md');
            const secondIsFile = secondEl.includes('.md');
            if ((firstIsFile && secondIsFile) || (!firstIsFile && !secondIsFile)) {
                return 0;
            }
            if (firstIsFile) {
                return -1;
            }
            if (secondIsFile) {
                return 1;
            }
            return undefined;
        });
        for (const path of paths) {
            if (path.includes('.md', 1)) {
                if (!path.includes('README.md', 1) && !path.includes('all.md', 1) && !path.includes('/modules/', 1) && !path.includes('globals.md', 1)) {
                    // Read file content and concat to new file
                    const content = fs.readFileSync(path);
                    fs.appendFileSync(allMdPath, content);
                    fs.appendFileSync(allMdPath, '\n');
                    fs.appendFileSync(allMdPath, '\n');
                    fs.appendFileSync(allMdPath, '<hr />');
                    fs.appendFileSync(allMdPath, '\n');
                    fs.appendFileSync(allMdPath, '\n');
                }
                if (!path.includes('README.md', 1)) {
                    fs.unlinkSync(path);
                }
            } else {
                fs.rmdirSync(path);
            }
        }

        // Find/replace relative links with hash links
        const docsBuff = fs.readFileSync(allMdPath);
        let docs = docsBuff.toString();
        docs = docs.replace(/\]\(((?!.*(github.com)).*)(#.*\))/g, ']($3');
        docs = docs.replace(/\]\(..\/interfaces\/.*?\.(.*?)\.md\)/g, '](#interface-$1)');
        docs = docs.replace(/\]\(..\/classes\/.*?\.(.*?)\.md\)/g, '](#class-$1)');
        docs = docs.replace(/\]\(..\/enums\/.*?\.(.*?)\.md\)/g, '](#enumeration-$1)');
        docs = docs.replace(/\]\(_types_\.(.*?)\.md\)/g, '](#interface-$1)');
        docs = docs.replace(/\]\(.*\.(.*?)\.md\)/g, '](#class-$1)');
        fs.writeFileSync(allMdPath, docs);
    });
}

main();
