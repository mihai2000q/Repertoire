import { forwardRef, SVGProps } from 'react'

interface CustomIconRhoadsGuitarProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconRhoadsGuitar = forwardRef<SVGSVGElement, CustomIconRhoadsGuitarProps>(
  ({ color, size = 24, strokeWidth = 2, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink={'http://www.w3.org/1999/xlink'}
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      strokeLinecap="round"
      strokeLinejoin="round"
      viewBox="0 0 512 512"
      xmlSpace="preserve"
      {...props}
    >
      <g>
        <g>
          <path
            d="M505.386,2.661c-1.794-2.403-4.978-3.3-7.764-2.189L397.183,40.526c-1.779,0.709-3.043,2.314-3.318,4.209
			c-0.274,1.896,0.485,3.792,1.989,4.977l8.954,7.048L284.187,210.004L14.604,356.06c-6.455,3.497-10.095,10.602-9.163,17.883
			c0.932,7.281,6.243,13.241,13.371,15l109.601,27.054c29.949,7.392,53.661,30.23,62.175,59.88l6.044,21.051
			c2.584,8.999,10.857,15.167,20.219,15.071c9.362-0.096,17.509-6.428,19.91-15.477l68.747-259.071L431.311,77.621l3.682,2.898
			c5.237,4.122,12.824,3.218,16.945-2.018l53.346-67.774C507.139,8.371,507.181,5.064,505.386,2.661z M212.075,441.106
			c-9.314,0-16.864-7.551-16.864-16.864s7.551-16.864,16.864-16.864s16.864,7.551,16.864,16.864
			C228.939,433.557,221.389,441.106,212.075,441.106z M231.788,385.503c-5.318,6.756-15.107,7.925-21.866,2.605l-55.921-44.016
			c-6.757-5.319-7.924-15.108-2.604-21.866c5.317-6.757,15.107-7.925,21.866-2.605l55.921,44.016
			C235.943,368.956,237.109,378.746,231.788,385.503z"
          />
        </g>
      </g>
    </svg>
  )
)

CustomIconRhoadsGuitar.displayName = 'CustomIconRhoadsGuitar'

export default CustomIconRhoadsGuitar
