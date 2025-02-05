import { mantineRender } from '../../../test-utils.tsx'
import SettingsModalCustomizationTab from './SettingsModalCustomizationTab.tsx'
import { screen } from '@testing-library/react'

describe('Settings Modal Customization Tab', () => {
  it('should render', () => {
    mantineRender(<SettingsModalCustomizationTab />)

    expect(screen.getByText(/feature coming soon/i))
  })
})
