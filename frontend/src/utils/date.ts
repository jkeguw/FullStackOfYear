export const formatDate = (date: Date | string) => {
  const d = new Date(date)
  return d.toLocaleDateString()
}

export const formatDateTime = (date: Date | string) => {
  const d = new Date(date)
  return d.toLocaleString()
}