# Cognito Hosted UI Customization

This directory contains the custom styling and branding assets for the SARC-NG Cognito Hosted Login UI.

## Files

- **`custom.css`** - Custom CSS stylesheet for the Cognito Hosted UI
- **`logo.svg`** - SARC-NG logo (SVG format recommended for best quality)
- **`upload-ui-assets.sh`** - Script to upload assets to S3 and generate pre-signed URLs
- **`apply-ui-customization.sh`** - Script to apply UI customization to Cognito User Pool

## Quick Start

### 1. Upload Assets to S3

First, upload your custom CSS and logo to an S3 bucket:

```bash
cd /workspaces/sarc-ng/infrastructure/sam/ui-customization
chmod +x upload-ui-assets.sh
./upload-ui-assets.sh dev sarc-ng-ui-assets-dev
```

This will:
- Create an S3 bucket (if it doesn't exist)
- Upload the CSS and logo files
- Generate pre-signed URLs valid for 7 days
- Save URLs to `ui-assets-urls-{env}.txt`

### 2. Apply UI Customization

After deploying your Cognito stack, apply the UI customization:

```bash
chmod +x apply-ui-customization.sh
./apply-ui-customization.sh dev
```

Or manually using AWS CLI:

```bash
# Get User Pool ID
USER_POOL_ID=$(aws cloudformation describe-stacks \
  --stack-name sarc-ng-cognito-dev \
  --query 'Stacks[0].Outputs[?OutputKey==`CognitoUserPoolId`].OutputValue' \
  --output text)

# Get Client ID
CLIENT_ID=$(aws cloudformation describe-stacks \
  --stack-name sarc-ng-cognito-dev \
  --query 'Stacks[0].Outputs[?OutputKey==`CognitoUserPoolClientId`].OutputValue' \
  --output text)

# Apply customization (CSS only)
aws cognito-idp set-ui-customization \
  --user-pool-id "${USER_POOL_ID}" \
  --client-id "${CLIENT_ID}" \
  --css "$(cat custom.css)"

# Apply customization (with logo)
aws cognito-idp set-ui-customization \
  --user-pool-id "${USER_POOL_ID}" \
  --client-id "${CLIENT_ID}" \
  --css "$(cat custom.css)" \
  --image-file fileb://logo.svg
```

## Customization Options

### CSS Customization

The `custom.css` file includes customizable CSS variables at the top:

```css
:root {
  --primary-color: #0066cc;      /* Main brand color */
  --primary-hover: #0052a3;      /* Hover state */
  --secondary-color: #6c757d;    /* Secondary elements */
  --success-color: #28a745;      /* Success messages */
  --danger-color: #dc3545;       /* Error messages */
  --background-color: #f8f9fa;   /* Page background */
  --card-background: #ffffff;    /* Form card background */
  --text-primary: #212529;       /* Primary text */
  --text-secondary: #6c757d;     /* Secondary text */
  --border-color: #dee2e6;       /* Borders */
  --border-radius: 8px;          /* Rounded corners */
  --box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
```

Update these values to match your brand colors.

### Logo Customization

Replace `logo.svg` with your own logo. Recommended specifications:
- **Format**: SVG (preferred) or PNG
- **Max size**: 2 MB
- **Recommended dimensions**: 200x80 pixels (width x height)
- **Background**: Transparent

### Background Gradient

The login page uses a purple gradient by default. To change it, edit the body background in `custom.css`:

```css
body {
  background: linear-gradient(135deg, #YOUR_COLOR_1 0%, #YOUR_COLOR_2 100%);
}
```

## Features

The custom UI includes:

✅ **Modern Design** - Clean, professional look with smooth animations
✅ **Responsive** - Works on mobile, tablet, and desktop
✅ **Accessible** - WCAG compliant with keyboard navigation support
✅ **Dark Mode** - Automatic dark mode support (respects system preference)
✅ **Custom Branding** - Logo, colors, and typography
✅ **Enhanced UX** - Better error messages, loading states, and transitions
✅ **Social Login Ready** - Styled buttons for social identity providers
✅ **MFA Support** - Custom styling for multi-factor authentication

## Testing the UI

After applying customization, test the hosted UI:

```bash
# Get the Hosted UI URL
HOSTED_UI_URL=$(aws cloudformation describe-stacks \
  --stack-name sarc-ng-cognito-dev \
  --query 'Stacks[0].Outputs[?OutputKey==`CognitoHostedUIUrl`].OutputValue' \
  --output text)

echo "Open this URL in your browser:"
echo "${HOSTED_UI_URL}"
```

## Updating the UI

To update the customization:

1. Edit `custom.css` or `logo.svg`
2. Re-run the apply script:
   ```bash
   ./apply-ui-customization.sh dev
   ```
3. Hard refresh your browser (Ctrl+Shift+R or Cmd+Shift+R)

## Limitations

- CSS file must be less than 100 KB
- Logo must be less than 2 MB
- Only specific CSS classes can be customized (see AWS documentation)
- Changes may take a few minutes to propagate

## Cognito CSS Classes

The following CSS classes can be customized:

| Class Name | Description |
|-----------|-------------|
| `.banner-customizable` | Top banner area |
| `.logo-customizable` | Logo image |
| `.label-customizable` | Form labels |
| `.inputField-customizable` | Input fields |
| `.button-customizable` | All buttons |
| `.submitButton-customizable` | Submit button specifically |
| `.secondaryButton-customizable` | Secondary actions |
| `.errorMessage-customizable` | Error messages |
| `.successMessage-customizable` | Success messages |
| `.redirect-customizable` | Links |
| `.socialButton-customizable` | Social login buttons |

## Resources

- [AWS Cognito UI Customization Documentation](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-ui-customization.html)
- [Cognito Hosted UI CSS Classes Reference](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-ui-customization.html#cognito-user-pools-app-ui-customization-css)

## Troubleshooting

### CSS not applying
- Check CSS file size (must be < 100 KB)
- Verify CSS syntax is valid
- Clear browser cache and hard refresh

### Logo not showing
- Check file size (must be < 2 MB)
- Verify file format (SVG, PNG, or JPG)
- Ensure logo has proper dimensions

### Changes not visible
- Wait 2-5 minutes for CDN propagation
- Clear browser cache (Ctrl+Shift+Delete)
- Try incognito/private mode

### Error: "InvalidParameterException"
- Check User Pool ID and Client ID are correct
- Ensure you have proper IAM permissions
- Verify the stack has been deployed

## Examples

See the included `custom.css` for a complete example with:
- Modern gradient background
- Rounded corners and shadows
- Smooth hover animations
- Responsive design
- Dark mode support
- Comprehensive form styling
