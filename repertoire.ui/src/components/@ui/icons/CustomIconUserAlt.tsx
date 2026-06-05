import { SVGProps } from 'react'

interface CustomIconUserAltProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

function CustomIconUserAlt({
  color,
  size = 24,
  strokeWidth = 2,
  ...props
}: CustomIconUserAltProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      strokeWidth={0.5}
      viewBox="0 0 24 24"
      xmlSpace={'preserve'}
      {...props}
    >
      <path d="M7.25 6a4.75 4.75 0 1 1 9.5 0 4.75 4.75 0 0 1-9.5 0ZM2.25 22c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75a.75.75 0 0 1-.75.75H3a.75.75 0 0 1-.75-.75Z" />
    </svg>
  )
}

export default CustomIconUserAlt
