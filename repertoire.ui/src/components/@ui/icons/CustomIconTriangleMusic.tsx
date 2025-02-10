import { forwardRef, SVGProps } from 'react'

interface CustomIconTriangleMusicProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconTriangleMusic = forwardRef<SVGSVGElement, CustomIconTriangleMusicProps>(
  ({ color, size = 24, strokeWidth = 12, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink={"http://www.w3.org/1999/xlink"}
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      strokeWidth={strokeWidth}
      version={'1.1'}
      viewBox="0 0 512 512"
      xmlSpace={'preserve'}
      {...props}
    >
      <g>
        <g>
          <path d="M445.969,440.543l-92.753-163.821l-25.205,25.205l73.972,130.651H46.639l177.672-313.807l51.014,90.101l25.205-25.205
			L239.258,75.449c-3.049-5.384-8.758-8.713-14.947-8.713c-6.189,0-11.897,3.329-14.947,8.713L2.229,441.291
			c-3.011,5.317-2.969,11.835,0.108,17.113c3.078,5.279,8.728,8.525,14.839,8.525h414.271c0.007,0,0.016,0,0.022,0
			c9.487,0,17.176-7.69,17.176-17.176C448.645,446.364,447.665,443.205,445.969,440.543z"/>
        </g>
      </g>
      <g>
        <g>
          <path d="M506.969,50.1c-6.707-6.707-17.583-6.707-24.29,0L258.102,274.676c-13.5-4.548-29.007-1.452-39.763,9.304
			c-15.118,15.118-15.118,39.631,0,54.749s39.631,15.118,54.749,0c10.755-10.755,13.851-26.263,9.304-39.763L506.969,74.389
			C513.677,67.682,513.677,56.807,506.969,50.1z"/>
        </g>
      </g>    </svg>
  )
)

CustomIconTriangleMusic.displayName = 'CustomIconTriangleMusic'

export default CustomIconTriangleMusic
