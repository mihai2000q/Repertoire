import { forwardRef, SVGProps } from 'react'

interface CustomIconDrumsProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconDrums = forwardRef<SVGSVGElement, CustomIconDrumsProps>(
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
      viewBox="0 0 75.904 48.843"
      xmlSpace={'preserve'}
      {...props}
    >
      <path
        d="M74.52,8.914H64.438V8.227H74.52c0.765,0,1.384-0.617,1.384-1.378c0-0.761-0.62-1.378-1.384-1.378H64.438V1.335
	C64.438,0.598,63.837,0,63.097,0c-0.74,0-1.342,0.598-1.342,1.335V5.47H51.675c-0.764,0-1.383,0.617-1.383,1.378
	c0,0.762,0.619,1.378,1.383,1.378h10.081v0.687H51.675c-0.764,0-1.383,0.619-1.383,1.38c0,0.762,0.619,1.376,1.383,1.376h10.081
	v26.426c-0.104,0.058-0.21,0.109-0.299,0.197l-3.896,3.881l-1.419-1.412c2.812-3.162,4.533-7.307,4.533-11.861
	c0-9.896-8.059-17.917-17.999-17.917c-6.731,0-12.591,3.686-15.678,9.136v-0.868c0-1.523-1.24-2.754-2.77-2.754H2.769
	C1.24,16.498,0,17.729,0,19.252v2.756c0,1.523,1.24,2.757,2.77,2.757h9.388V38.13c-0.104,0.059-0.211,0.109-0.299,0.199
	l-6.854,6.821c-0.523,0.521-0.523,1.367,0,1.888C5.266,47.3,5.608,47.43,5.952,47.43c0.343,0,0.688-0.13,0.949-0.392l5.256-5.231
	v5.7c0,0.736,0.602,1.336,1.342,1.336s1.341-0.6,1.341-1.336v-5.701l5.258,5.232c0.262,0.262,0.604,0.392,0.947,0.392
	c0.345,0,0.687-0.13,0.948-0.392c0.524-0.521,0.524-1.367,0-1.888l-6.852-6.821c-0.089-0.088-0.196-0.142-0.302-0.199V24.765h9.39
	c0.355,0,0.69-0.08,1.004-0.201c-0.35,1.391-0.554,2.838-0.554,4.336c0,4.555,1.722,8.7,4.533,11.861l-4.408,4.389
	c-0.524,0.521-0.524,1.367,0,1.889c0.262,0.26,0.604,0.391,0.947,0.391c0.345,0,0.688-0.131,0.949-0.391l4.425-4.406
	c3.125,2.609,7.152,4.186,11.55,4.186c4.399,0,8.426-1.576,11.553-4.186l1.435,1.429l-1.058,1.054c-0.524,0.521-0.524,1.367,0,1.889
	c0.262,0.261,0.604,0.391,0.948,0.391c0.344,0,0.687-0.13,0.947-0.391l1.06-1.053l1.094,1.088c0.262,0.26,0.604,0.391,0.948,0.391
	c0.344,0,0.687-0.131,0.949-0.391c0.522-0.521,0.522-1.367,0-1.889l-1.096-1.089l2.3-2.29v4.287c0,0.736,0.602,1.335,1.342,1.335
	c0.739,0,1.341-0.599,1.341-1.335v-4.289l5.256,5.232c0.262,0.262,0.604,0.392,0.948,0.392c0.344,0,0.687-0.13,0.947-0.392
	c0.523-0.521,0.523-1.367,0-1.887l-6.853-6.822c-0.089-0.09-0.196-0.142-0.3-0.199V11.67H74.52c0.765,0,1.384-0.614,1.384-1.376
	C75.904,9.533,75.284,8.914,74.52,8.914z"
      />
    </svg>
  )
)

CustomIconDrums.displayName = 'CustomIconDrums'

export default CustomIconDrums
