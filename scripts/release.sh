#!/bin/bash
set -e

VERSION=${1:-$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)}
DIST_DIR="dist"

echo "Creating release artifacts for version ${VERSION}"

# Create dist directory
mkdir -p "${DIST_DIR}"

# Package each platform/arch combination
for os in darwin linux windows; do
    for arch in amd64 arm64; do
        if [ "$os" == "windows" ]; then
            binary_name="helm-cm-push.exe"
        else
            binary_name="helm-cm-push"
        fi

        binary_path="bin/${os}/${arch}/helm-cm-push"

        # Check if binary exists
        if [ ! -f "${binary_path}" ]; then
            echo "Skipping ${os}/${arch} - binary not found"
            continue
        fi

        archive_name="helm-push_${VERSION}_${os}_${arch}.tar.gz"
        temp_dir=$(mktemp -d)

        echo "Creating ${archive_name}..."

        # Create proper structure
        mkdir -p "${temp_dir}/bin"
        cp "${binary_path}" "${temp_dir}/bin/${binary_name}"
        cp LICENSE "${temp_dir}/"
        cp plugin.yaml "${temp_dir}/"

        # Create tarball
        tar -czf "${DIST_DIR}/${archive_name}" -C "${temp_dir}" .

        # Cleanup
        rm -rf "${temp_dir}"

        echo "✓ ${archive_name}"
    done
done

echo ""
echo "Release artifacts created in ${DIST_DIR}:"
ls -lh "${DIST_DIR}/"
