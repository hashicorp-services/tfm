#!/usr/bin/env bash
set -e

echo "Step 1: Remove unfiltered subtree..."
if [ -d "./github/skills/go-development" ]; then
  git rm -r .github/skills/go-development
  git commit -m "remove: delete unfiltered subtree"
  echo "✓ Flattened structure"
fi

echo "Step 2: Flatten the directory structure..."
mkdir ./tmp
cd ./tmp
git clone https://github.com/netresearch/go-development-skill.git go-dev-temp
cd go-dev-temp
cp LICENSE skills/go-development/LICENSE

# Filter to just the skills/go-development directory using git filter-branch
git filter-branch --subdirectory-filter skills/go-development -- --all

echo "Step 3: move structure..."
# Go back to tfm repo
cd /Users/abuxton/src/hashicorp-services/tfm
# Move nested go-development folder content up
if [ -d "./tmp/go-dev-temp/skills/go-development" ]; then
    mv ./tmp/go-dev-temp/skills/go-development .github/skills/go-development
    rmdir ./tmp/go-dev-temp/
    echo "✓ deployed skills/go-development to .github/skills/go-development"
else
    echo "✗ Error: ./tmp/go-dev-temp/skills/go-development not found"
    exit 1
fi


echo "Step 3: Stage and commit changes..."
git add .github/skills/go-development
git commit -m "chore: flatten go-development skill structure"

echo "Step 4: Push to remote..."
git push origin feature/specify-init

echo "✓ Complete!"




# Add the filtered version as a subtree
git subtree add --prefix .github/skills/go-development /tmp/go-dev-temp main --squash

# Clean up
rm -rf /tmp/go-dev-temp

git push origin feature/specify-init
