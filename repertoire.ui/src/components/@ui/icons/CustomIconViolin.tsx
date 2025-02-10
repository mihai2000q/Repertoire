import { forwardRef, SVGProps } from 'react'

interface CustomIconViolinProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
  strokeWidth?: number
}

const CustomIconViolin = forwardRef<SVGSVGElement, CustomIconViolinProps>(
  ({ color, size = 24, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink={'http://www.w3.org/1999/xlink'}
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      version="1.1"
      viewBox="0 0 512.001 512.001"
      xmlSpace={'preserve'}
      {...props}
    >
      <g>
        <g>
          <path
            d="M509.475,30.246l-27.72-27.72c-3.747-3.745-9.962-3.26-13.08,1.027L448.49,31.306l-3.48-3.48
			c-6.558-6.559-17.194-6.559-23.753,0s-6.559,17.193,0,23.753l5.726,5.726l-12.887,12.887l-5.726-5.726
			c-6.558-6.559-17.194-6.559-23.753,0s-6.559,17.193,0,23.753l5.726,5.726l-68.355,68.355
			c-36.895-25.628-84.727-24.823-113.873,4.323l-0.001,0.002c-9.224,9.224-18.915,17.859-29.028,25.911
			c7.19,18.739,1.137,43.33-16.863,61.329c-22.95,22.95-50.364,32.744-68.949,14.16c-3.894-3.894-6.805-8.458-8.781-13.441
			c-10.69,3.813-21.635,7.171-32.828,10.052c-11.683,3.006-21.983,8.573-30.257,16.846c-34.97,34.97-26.884,100.769,22.703,158.821
			L8.561,489.745c-2.385,3.328-2.523,7.948,0,11.456c3.163,4.399,9.295,5.403,13.695,2.238l49.442-35.548
			c58.051,49.585,123.851,57.672,158.821,22.703c8.275-8.275,13.84-18.574,16.846-30.257c2.87-11.153,6.243-22.099,10.082-32.816
			c-4.994-1.976-9.568-4.891-13.471-8.793c-18.584-18.584-8.79-45.999,14.159-68.949c18.013-18.013,42.627-24.061,61.37-16.846
			c8.061-10.159,16.687-19.863,25.871-29.047c29.146-29.147,29.951-76.977,4.323-113.872l68.355-68.355l5.726,5.726
			c6.558,6.559,17.193,6.56,23.753,0c6.559-6.559,6.559-17.193,0-23.753l-5.726-5.726l12.887-12.887l5.726,5.726
			c6.558,6.559,17.194,6.56,23.753,0c6.559-6.559,6.559-17.193,0-23.753l-3.48-3.48l27.754-20.184
			C512.734,40.209,513.223,33.994,509.475,30.246z M156.838,392.161c-0.197,2.516-1.486,4.818-3.527,6.303l-44.851,32.619
			c-3.434,2.497-8.174,2.125-11.176-0.877l-15.488-15.488c-3.002-3.002-3.374-7.743-0.877-11.176l32.619-44.851
			c1.485-2.041,3.788-3.33,6.303-3.527c2.516-0.197,4.992,0.717,6.777,2.5l27.721,27.721
			C156.121,387.17,157.035,389.645,156.838,392.161z M184.436,298.954c-12.93,0.945-22.166,6.451-39.408,12.45
			c-13.638,4.694-31.386,9.164-44.551,3.22c-4.457-1.773-8.869-5.555-11.027-10.584c3.097,0.868,8.017,2.4,12.816,1.462
			c11.862-1.845,20.771-6.791,36.446-12.233c13.669-4.695,31.316-9.136,44.546-3.197c4.558,1.81,8.398,5.209,11.032,10.56
			C193.191,300.327,188.058,298.561,184.436,298.954z M218.731,373.289c-5.442,15.675-10.388,24.584-12.233,36.446
			c-0.937,4.799,0.595,9.718,1.462,12.815c-5.029-2.157-8.811-6.571-10.584-11.027c-5.943-13.165-1.475-30.912,3.22-44.551
			c5.998-17.242,11.504-26.478,12.45-39.408c0.393-3.621-1.372-8.754-1.68-9.854c5.351,2.634,8.751,6.473,10.56,11.032
			C227.867,341.974,223.426,359.621,218.731,373.289z"
          />
        </g>
      </g>
    </svg>
  )
)

CustomIconViolin.displayName = 'CustomIconViolin'

export default CustomIconViolin
