#!/bin/bash
# parameters
files=1
warm=2

# cleanup
cleanup() {
    echo "cleaning up"
    # rm -rf /tmp/bench
}
trap cleanup EXIT

# check if hyperfine is installed
if ! command -v hyperfine &>/dev/null; then
    echo "hyperfine is not installed. Please install hyperfine to continue."
    exit 1
fi

# cloning candidates
echo ""
echo "clone SSGs"
echo ""
git clone --depth=1 https://github.com/shreyaskaundinya/garlic /tmp/bench/garlic

echo ""
echo "build SSGs"
echo ""

cd /tmp/bench/garlic && go build -o garlic && mv garlic /tmp/bench/garlic/garlic


# create the content folder
mkdir -p /tmp/bench/garlic/src/content

ls -la /tmp/bench/garlic/src/

ls -la /tmp/bench/garlic/src/content/
ls -la /tmp/bench/garlic/src/templates/
ls -la /tmp/bench/garlic/src/assets/
ls -la /tmp/bench/garlic/src/components/

# clean content/* dirs
echo ""
echo "Cleaning content directories"
echo ""
rm -rf /tmp/bench/garlic/src/content/*

# create multiple copies of the test file
echo ""
echo "Spawning $files different markdown files..."
for ((i = 0; i < files; i++)); do
    cp /tmp/bench/garlic/test/bench/bench.md "/tmp/bench/garlic/src/content/test_$i.md"
done

# run hyperfine
echo ""
echo "running benchmark: $files md files and $warm warmup runs"
echo ""

# "cd /tmp/bench/sapling/benchmark && ./../sapling run" \

hyperfine --ignore-failure -p 'sync' -w $warm \
  "cd /tmp/bench/garlic && ./garlic --src-folder ./src --dest-folder ./dest --seed-files"
echo ""