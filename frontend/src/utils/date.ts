export const formatDate = (date: Date | string) => {
  const d = new Date(date)
  return d.toLocaleDateString()
}