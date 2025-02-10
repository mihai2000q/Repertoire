import { forwardRef, SVGProps } from 'react'

interface CustomIconQuestionMarkProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconQuestionMark = forwardRef<SVGSVGElement, CustomIconQuestionMarkProps>(
  ({ color, size = 24, ...props }, ref) => (
    <svg
      ref={ref}
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 48 48"
      xmlSpace={'preserve'}
      {...props}
    >
      <g>
        <path d="M24.3,6A11.2,11.2,0,0,0,16,9.3a11,11,0,0,0-3.5,8.2,2.5,2.5,0,0,0,5,0,6.5,6.5,0,0,1,2-4.7A6.2,6.2,0,0,1,24.2,11a6.5,6.5,0,0,1,1,12.9,4.4,4.4,0,0,0-3.7,4.4v3.2a2.5,2.5,0,0,0,5,0V28.7a11.6,11.6,0,0,0,9-11.5A11.7,11.7,0,0,0,24.3,6Z" />
        <circle cx="24" cy="39.5" r="2.5" />
      </g>
    </svg>
  )
)

CustomIconQuestionMark.displayName = 'CustomIconQuestionMark'

export default CustomIconQuestionMark
