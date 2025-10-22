#!/bin/bash
# parameters
files=1000
warm=10

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

# check if hugo is installed
if ! command -v hugo &>/dev/null; then
    echo "hugo is not installed. Please install hugo to continue."
fi

# cloning candidates
echo ""
echo "clone SSGs"
echo ""
git clone --depth=1 https://github.com/anna-ssg/anna /tmp/bench/anna
git clone --depth=1 https://github.com/anirudhRowjee/saaru /tmp/bench/saaru
# git clone --depth=1 https://github.com/NavinShrinivas/sapling /tmp/bench/sapling
git clone --depth=1 https://github.com/shreyaskaundinya/garlic /tmp/bench/garlic

# copy benchmark file
cp /tmp/bench/anna/site/content/posts/bench.md /tmp/bench/test.md

echo ""
echo "build SSGs"
echo ""

cd /tmp/bench/anna && go build && cd /tmp/bench

# build rust based SSGs (edit this block if they are already installed)
# cd /tmp/bench/sapling && cargo build --release && mv target/release/sapling .
cd /tmp/bench/saaru && cargo build --release && mv target/release/saaru .

cd /tmp/bench/garlic && go build -o garlic && mv garlic /tmp/bench/garlic/garlic

## setup hugo
hugo new site /tmp/bench/hugo; cd /tmp/bench/hugo
hugo new theme mytheme; echo "theme = 'mytheme'" >> hugo.toml; cd /tmp/bench

## setup 11ty


# create the content folder for garlic
mkdir -p /tmp/bench/garlic/src/content

# clean content/* dirs
echo ""
echo "Cleaning content directories"
echo ""
rm -rf /tmp/bench/anna/site/content/posts/*
rm -rf /tmp/bench/saaru/docs/src/*
# rm -rf /tmp/bench/sapling/benchmark/content/blog/*
rm -rf /tmp/bench/hugo/content/*
rm -rf /tmp/bench/garlic/src/content/*

# create multiple copies of the test file
echo ""
echo "Spawning $files different markdown files..."
for ((i = 0; i < files; i++)); do
    cp /tmp/bench/test.md "/tmp/bench/anna/site/content/posts/test_$i.md"
    cp /tmp/bench/test.md "/tmp/bench/saaru/docs/src/test_$i.md"
    # cp /tmp/bench/test.md "/tmp/bench/sapling/benchmark/content/blogs/test_$i.md"
    cp /tmp/bench/test.md "/tmp/bench/hugo/content/test_$i.md"
    cp /tmp/bench/garlic/test/bench/bench.md "/tmp/bench/garlic/src/content/test_$i.md"
done

# run hyperfine
echo ""
echo "running benchmark: $files md files and $warm warmup runs"
echo ""

# "cd /tmp/bench/sapling/benchmark && ./../sapling run" \

hyperfine --ignore-failure -p 'sync' -w $warm \
  "cd /tmp/bench/hugo && hugo" \
  "cd /tmp/bench/anna && ./anna -r \"site/\"" \
  "cd /tmp/bench/saaru && ./saaru --base-path ./docs" \
  "cd /tmp/bench/garlic && ./garlic --src-folder /tmp/bench/garlic/src --dest-folder /tmp/bench/garlic/dest --seed-files"

echo ""