import fetch from "node-fetch";
const url = "http://localhost:8500";
async function main() {
  try {
    const response = await fetch(url);
    console.log({ response });
  } catch (error) {
    console.log(error);
  }
}

setInterval(() => {
  Promise.resolve(main());
}, 1_000);
