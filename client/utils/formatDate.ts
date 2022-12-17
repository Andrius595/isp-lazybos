export default (date: string) => {
  return `${date.split('T')[0]} ${date.split('T')[1].substring(0,8)}`
}