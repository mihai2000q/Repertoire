import { forwardRef, SVGProps } from 'react'

interface CustomIconArrowRightProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number
}

const CustomIconArrowRight = forwardRef<SVGSVGElement, CustomIconArrowRightProps>(
  ({ color, size = 24, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      strokeLinecap="round"
      strokeLinejoin="round"
      viewBox="0 0 1024 1024"
      {...props}
    >
      <path d="M419.3 264.8l-61.8 61.8L542.9 512 357.5 697.4l61.8 61.8L666.5 512z" />
    </svg>
  )
)

CustomIconArrowRight.displayName = 'CustomIconArrowRight'

export default CustomIconArrowRight
