import { mantineRender } from '../../../test-utils.tsx'
import SelectableAvatar from './SelectableAvatar.tsx'
import { screen } from '@testing-library/react'

describe('Selectable Avatar', () => {
  it('should not display checkmark when not selected', () => {
    const imgAlt = 'something'

    mantineRender(<SelectableAvatar isSelected={false} alt={imgAlt} src={'something.png'} />)

    expect(screen.getByTestId('selected-checkmark')).not.toBeVisible()
    expect(screen.getByRole('img', { name: imgAlt })).toBeInTheDocument()
  })

  it('should display checkmark when selected', async () => {
    const imgAlt = 'something'

    mantineRender(<SelectableAvatar isSelected={true} alt={imgAlt} src={'something.png'} />)

    expect(screen.getByTestId('selected-checkmark')).toBeVisible()
    expect(screen.getByRole('img', { name: imgAlt })).toBeInTheDocument()
  })
})
