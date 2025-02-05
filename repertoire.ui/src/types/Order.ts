import { ReactElement } from 'react'

export default interface Order {
  label: string
  value: string
  icon?: ReactElement
  property?: string
}
