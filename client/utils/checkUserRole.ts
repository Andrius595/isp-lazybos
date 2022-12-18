export default (user, role) => {
  if (!user) {
    return false
  }

  if (!user.role) {
    return role === 'user'
  }

  return role === user.role
}