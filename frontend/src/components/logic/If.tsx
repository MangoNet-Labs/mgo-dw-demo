import { ReactNode } from 'react'

interface Props {
  prediction: boolean
  children: ReactNode | ReactNode[]
}

export default function If({ prediction, children }: Props) {
  if (!prediction) return null
  return children
}
