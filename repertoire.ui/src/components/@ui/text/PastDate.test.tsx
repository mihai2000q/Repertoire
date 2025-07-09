import dayjs from 'dayjs'
import PastDate from './PastDate'
import { mantineRender } from '../../../test-utils.tsx'

describe('Past Date', () => {
  const mockNow = dayjs('2024-01-17T12:00:00') // Wednesday, January 17, 2024

  beforeEach(() => {
    vi.useFakeTimers()
    vi.setSystemTime(mockNow.toDate())
  })

  afterEach(() => vi.useRealTimers())

  it('should render "Today" for current date', () => {
    const today = mockNow.toISOString()
    const { getByText } = mantineRender(<PastDate dateValue={today} />)
    expect(getByText('Today')).toBeInTheDocument()
  })

  it('should render "Yesterday" for previous day', () => {
    const yesterday = mockNow.subtract(1, 'day').toISOString()
    const { getByText } = mantineRender(<PastDate dateValue={yesterday} />)
    expect(getByText('Yesterday')).toBeInTheDocument()
  })

  it('should render weekday name for dates in the same week', () => {
    const monday = mockNow.subtract(2, 'day').toISOString() // Sunday
    const { getByText } = mantineRender(<PastDate dateValue={monday} />)
    expect(getByText('Monday')).toBeInTheDocument()
  })

  it('should render formatted date for older dates', () => {
    const monday = mockNow.subtract(3, 'day').toISOString() // Sunday
    const { getByText } = mantineRender(<PastDate dateValue={monday} />)
    expect(getByText('14 Jan')).toBeInTheDocument()
  })

  it('should render formatted date for older dates even from last year', () => {
    const oldDate = '2023-12-25T00:00:00' // Christmas
    const { getByText } = mantineRender(<PastDate dateValue={oldDate} />)
    expect(getByText('25 Dec')).toBeInTheDocument()
  })

  it('uses custom date format when provided', () => {
    const oldDate = '2023-12-25T00:00:00'
    const { getByText } = mantineRender(<PastDate dateValue={oldDate} dateFormat="MMMM D, YYYY" />)
    expect(getByText('December 25, 2023')).toBeInTheDocument()
  })
})
