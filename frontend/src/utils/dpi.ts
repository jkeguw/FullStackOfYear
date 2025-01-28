export const convertDPI = (sourceDPI: number, targetDPI: number, sensitivity: number) => {
  return (sensitivity * sourceDPI) / targetDPI
}