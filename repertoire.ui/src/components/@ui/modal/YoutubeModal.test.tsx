import { reduxRender } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import YoutubeModal from './YoutubeModal.tsx'

describe('Youtube Modal', () => {
  // eslint-disable-next-line no-console
  const originalError = console.error

  beforeAll(() => {
    // eslint-disable-next-line no-console
    console.error = (...args: { toString: () => string | string[] }[]) => {
      if (
        args[0]?.toString().includes('AbortError') ||
        args[0]?.toString().includes('NetworkError')
      ) {
        return // Silently ignore
      }
      originalError(...args)
    }
  })

  afterAll(() => {
    // eslint-disable-next-line no-console
    console.error = originalError // Restore
  })

  it.each([
    [
      'https://www.youtube.com/watch?v=tAGnKpE4NCI',
      'https://www.youtube-nocookie.com/embed/tAGnKpE4NCI'
    ],
    ['https://youtu.be/tAGnKpE4NCI', 'https://www.youtube-nocookie.com/embed/tAGnKpE4NCI'],
    ['https://www.youtu.be/tAGnKpE4NCI', 'https://www.youtube-nocookie.com/embed/tAGnKpE4NCI']
  ])('should render', async (link, expectedLinkToRender) => {
    const title = 'Some Song'

    reduxRender(<YoutubeModal title={title} link={link} opened={true} onClose={() => {}} />)

    expect(await screen.findByRole('dialog', { name: title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: title })).toBeInTheDocument()
    expect(screen.getByTitle('Embedded Youtube')).toHaveAttribute('src', expectedLinkToRender)
  })
})
