export type Role = 'engineer' | 'manager' | 'observer'

export type UserMe = {
  user_id: number
  role: Role
}

export type Project = {
  id: number
  name: string
  description: string
}

export type Defect = {
  id: number
  title: string
  description?: string
  project_id: number
  assigned_to?: number | null
  status: 'new' | 'in_progress' | 'on_review' | 'closed'
  priority: 'low' | 'medium' | 'high'
}


