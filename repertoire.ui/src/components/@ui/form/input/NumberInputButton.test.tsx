import { mantineRender } from '../../../../test-utils.tsx'
import NumberInputButton from './NumberInputButton.tsx'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Number Input Button', () => {
  it('should render', async () => {
    const user = userEvent.setup()

    const ariaLabel = 'number-input-label'
    const inputAriaLabel = 'number-input-input'

    const selectedTooltip = 'it is selected'
    const defaultTooltip = 'it is default'

    const { rerender } = mantineRender(
      <NumberInputButton
        aria-label={ariaLabel}
        inputProps={{ 'aria-label': inputAriaLabel }}
        tooltipLabels={{ selected: selectedTooltip, default: defaultTooltip }}
      />
    )

    const button = screen.getByRole('button', { name: ariaLabel })
    expect(button).toBeInTheDocument()

    await user.hover(button)
    await waitFor(() =>
      expect(screen.getByRole('tooltip', { name: defaultTooltip })).toBeInTheDocument()
    )

    await user.click(button)
    expect(screen.getByRole('textbox', { name: inputAriaLabel })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: inputAriaLabel })).toHaveFocus()

    rerender(
      <NumberInputButton
        aria-label={ariaLabel}
        inputProps={{ 'aria-label': inputAriaLabel }}
        tooltipLabels={{ selected: selectedTooltip, default: defaultTooltip }}
        isSelected={true}
      />
    )

    await user.hover(button)
    await waitFor(() =>
      expect(screen.getByRole('tooltip', { name: selectedTooltip })).toBeInTheDocument()
    )
  })
})
