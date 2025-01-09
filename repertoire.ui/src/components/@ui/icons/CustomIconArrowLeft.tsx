import { forwardRef, SVGProps } from 'react'

interface CustomIconArrowLeftProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number
}

const CustomIconArrowLeft = forwardRef<SVGSVGElement, CustomIconArrowLeftProps>(
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
      <path d="M604.7 759.2l61.8-61.8L481.1 512l185.4-185.4-61.8-61.8L357.5 512z" />
    </svg>
  )
)

CustomIconArrowLeft.displayName = "CustomIconArrowLeft";

export default CustomIconArrowLeft
