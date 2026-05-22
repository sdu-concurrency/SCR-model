import { readFileSync, writeFileSync } from 'node:fs';

const mainTranslate = JSON.parse(readFileSync(new URL('./main_translation.json', import.meta.url)));
const locales = ["da", "en"]

function* traverse(o, path = []) {
    for (var i of Object.keys(o)) {
        const itemPath = path.concat(i);
        yield [i, o[i], itemPath];
        if (o[i] !== null && typeof (o[i]) == "object") {
            yield* traverse(o[i], itemPath);
        }
    }
}

function assignNested(base, keys, value) {
    for (let i = 0; i < keys.length - 1; i++) {
        if (i === keys.length - 2) {
            base[keys[i]] = value
            break
        }
        base = base[keys[i]] = base[keys[i]] || {};
    }
}

const res = locales.reduce((prev, curr) => { prev[curr] = {}; return prev }, {})
for (var [key, value, path] of traverse(mainTranslate)) {
    if (typeof value === "string" && locales.includes(key)) {
        assignNested(res[key], path, value)
    }
}

for (const [key, value] of Object.entries(res)) {

    const file = new URL(`i18n/${key}.json`, import.meta.url)
    writeFileSync(file, JSON.stringify(value, null, 2))

}