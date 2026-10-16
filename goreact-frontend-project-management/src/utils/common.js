const common = {
  sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  },
  bytesToMb(bytes) {
    if (bytes === 0) {
      return '0 MB';
    }
    const megabytes = bytes / (1024 * 1024);
    return megabytes.toFixed(0) + ' MB';
  },
};

export default common;
