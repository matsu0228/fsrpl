import { execSync } from "child_process";
import * as firebase from "@firebase/rules-unit-testing";
import * as admin from "firebase-admin";

const emulatorHost = process.env.FIRESTORE_EMULATOR_HOST;
const projectId = "firestore-test"; // NOTE: set unique projectId in each tests
const fsrpl = "go run ../..."; // NOTE: set `fsrpl` commands on your environments
const restoreCmd = `FIRESTORE_EMULATOR_HOST=${emulatorHost} ${fsrpl} restore --path "./testData/" "test/*" --debug --emulators-project-id=${projectId}`;

// production code. For simplicity, wrote the same file with test code.
const calculate = async (
  firestore: admin.firestore.Firestore | any,
  category: string
): Promise<number> => {
  let total = 0;
  await firestore
    .collection("test")
    .where("category", "==", category)
    .get()
    .then((qs: any) =>
      qs.docs.forEach((doc: any) => {
        const product = doc.data();
        if (typeof product.price !== "number") {
          console.log("WARN: unexpected type: ", product.price);
          return;
        }
        total += product.price;
      })
    );
  return total;
};

describe("jest example", () => {
  const testFirestore = firebase.initializeAdminApp({ projectId }).firestore();
  beforeAll(async () => {
    const stdout = await execSync(`${restoreCmd}`);
    console.log(`imported test data: ${stdout.toString()}`);
  });

  afterAll(async () => {
    await firebase.clearFirestoreData({
      projectId,
    });
    console.log("cleared test data: ", projectId);
  });

  test("succeed to import of test data", async () => {
    const gotDocIds = await testFirestore
      .collection("test")
      .get()
      .then((qs) => qs.docs.map((doc) => doc.id));
    console.log("got documents id: ", gotDocIds);
    expect(gotDocIds).toEqual(expect.arrayContaining(["pc1", "pc2", "sp1"]));
  });

  test("calculate computer products price (10000 + 40000)", async () => {
    const expected = 50000;
    const gotTotal = await calculate(testFirestore, "computer");
    expect(gotTotal).toBe(expected);
  });
});
