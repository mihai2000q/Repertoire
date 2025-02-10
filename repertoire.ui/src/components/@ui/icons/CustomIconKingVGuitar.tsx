import { forwardRef, SVGProps } from 'react'

interface CustomIconKingVGuitarProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconKingVGuitar = forwardRef<SVGSVGElement, CustomIconKingVGuitarProps>(
  ({ color, size = 24, strokeWidth = 2, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink="http://www.w3.org/1999/xlink"
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      strokeLinecap="round"
      strokeLinejoin="round"
      viewBox="0 0 484.746 484.746"
      xmlSpace="preserve"
      {...props}
    >
      <g>
        <g>
          <path
            d="M16.321,300.903l121.595,7.927c10.788,0.705,20.091,10.015,20.773,20.807l9.349,146.476
			c0.689,10.78,4.921,11.577,9.454,1.763l109.653-237.396c4.528-9.814,8.308-17.853,8.384-18.001
			c0.071-0.15,0.071-0.307,0.084-0.336c0.016-0.032-0.021-0.096-0.068-0.167c-0.053-0.068-0.044-0.068,0.008-0.006
			c0.048,0.06-0.12-0.18-0.353-0.561c-0.236-0.375-0.461-0.782-0.521-0.886c-0.061-0.11-4.453-8.674-2.04-19.208
			c0.309-1.347,0.773-2.615,1.19-3.909l96.534-94.88c11.598-7.743,28.188-1.801,28.188-1.801
			c-2.677-19.278,50.714-77.403,57.182-83.877c6.48-6.478-10.623-7.169-23.074-9.542C440.197,4.945,440.385,0,440.385,0
			c-22.934,33.015-75.227,56.841-75.227,56.841c5.213,10.247,3.662,18.685,1.631,23.798l-90.147,88.614
			c-1.214,0.296-2.416,0.557-3.583,0.737c-9.758,1.563-17.824-0.385-17.868-0.361c-0.045,0.022-0.116,0.054-0.161,0.076
			c-0.048,0.022-7.896,3.983-17.552,8.839L14.273,290.841C4.618,295.697,5.535,300.198,16.321,300.903z"
          />
        </g>
      </g>
    </svg>
  )
)

CustomIconKingVGuitar.displayName = 'CustomIconKingVGuitar'

export default CustomIconKingVGuitar
