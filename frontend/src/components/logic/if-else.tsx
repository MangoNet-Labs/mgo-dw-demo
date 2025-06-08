import { ReactElement } from 'react'

interface Props {
  prediction: boolean
  children: [ReactElement, ReactElement]
}
export default function IfElse({ prediction, children }: Props) {
  return prediction ? children[0] : children[1]
}
