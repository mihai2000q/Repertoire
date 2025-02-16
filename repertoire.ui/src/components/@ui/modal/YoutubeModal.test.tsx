import { reduxRender } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import YoutubeModal from './YoutubeModal.tsx'

describe('Youtube Modal', () => {
  it('should render', async () => {
    const title = 'Some Song'
    const link = 'https://www.youtube.com/watch?v=tAGnKpE4NCI'
    const expectedLinkToRender = 'https://www.youtube.com/embed/tAGnKpE4NCI'

    reduxRender(<YoutubeModal title={title} link={link} opened={true} onClose={() => {}} />)

    expect(await screen.findByRole('dialog', { name: title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: title })).toBeInTheDocument()
    expect(screen.getByTitle('Embedded Youtube')).toHaveAttribute('src', expectedLinkToRender)
  })
})
