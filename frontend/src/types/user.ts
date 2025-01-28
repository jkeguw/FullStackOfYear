export interface User {
  id: string
  username: string
  role: string
}

export interface UserState {
  user: User | null
  token: string | null
}