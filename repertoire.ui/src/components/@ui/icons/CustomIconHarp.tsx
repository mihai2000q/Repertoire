import { forwardRef, SVGProps } from 'react'

interface CustomIconHarpProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconHarp = forwardRef<SVGSVGElement, CustomIconHarpProps>(
  ({ color, size = 24, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink={'http://www.w3.org/1999/xlink'}
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      version={'1.1'}
      viewBox="0 0 512 512"
      xmlSpace={'preserve'}
      {...props}
    >
      <g>
        <path
          d="M456.967,17.542c0,0-27.005-39.468-99.721,0c-21.912,11.9-110.116,69.414-192.168,70.648
		c-69.594,1.031-121.542-56.1-109.07-6.244c13.69,54.761,105.949,398.897,105.949,398.897h-22.852V512h112.192L112.101,125.589
		c0,0,34.068,6.379,89.333,2.069c66.48-5.198,120.496-49.864,162.048-33.24c41.553,16.624-20.775,220.225-126.732,334.478
		l12.472,37.4C249.221,466.296,473.591,283.47,456.967,17.542z"
        />
        <path d="M146.145,195.613l5.883,16.33l36.835-75.146c-4.738,0.248-9.455,0.444-14.126,0.519L146.145,195.613z" />
        <path d="M250.921,126.01c-5.266,1.625-10.621,3.13-16.082,4.468l-64.63,131.892l5.882,16.324L250.921,126.01z" />
        <path
          d="M319.063,102.744c-5.28,1.392-10.704,3.136-16.375,5.092h-0.016L194.242,329.092l5.882,16.338L319.063,102.744
		z"
        />
        <polygon points="367.784,119.202 367.784,119.18 361.843,102.978 218.283,395.85 224.165,412.158 	" />
      </g>
    </svg>
  )
)

CustomIconHarp.displayName = 'CustomIconHarp'

export default CustomIconHarp
