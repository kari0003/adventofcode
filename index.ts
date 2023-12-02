import { runner } from "./2023/2";
import { writeFileSync } from "fs";

runner.run();

const output = JSON.stringify(runner.context);
writeFileSync("./output.json", output);
console.log(output);
