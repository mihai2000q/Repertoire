import { forwardRef, SVGProps } from 'react'

interface CustomIconLightningTrioProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number
}

const CustomIconLightningTrio = forwardRef<SVGSVGElement, CustomIconLightningTrioProps>(
  ({ color, size = 24 }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      strokeLinecap="round"
      strokeLinejoin="round"
      viewBox="0 0 32 32"
    >
      <path d="M3.385 3.897l15.703 20.922-0.304-3.73 10.456 8.282-11.526-17.404 0.304 3.73-14.634-11.8zM4.843 3.041l16.717 7.313-1.073-2.295 8.437 0.865-11.293-6.35 1.073 2.295-13.861-1.829zM2.276 5.725l2.063 17.887 1.448-2.31 3.85 7.877-0.543-13.632-1.448 2.31-5.37-12.132z"></path>
    </svg>
  )
)

CustomIconLightningTrio.displayName = 'CustomIconLightningTrio'

export default CustomIconLightningTrio
