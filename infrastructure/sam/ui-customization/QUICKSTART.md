# SARC-NG Cognito Hosted UI Customization Guide

## 🎨 What You Get

A fully customized, modern login interface for your AWS Cognito authentication with:

- ✅ **Custom branding** with your logo and colors
- ✅ **Modern gradient background** (purple by default, easily customizable)
- ✅ **Responsive design** that works on mobile, tablet, and desktop
- ✅ **Professional animations** and smooth transitions
- ✅ **Dark mode support** (automatic based on system preference)
- ✅ **Enhanced UX** with better error messages and loading states
- ✅ **WCAG accessibility** compliant
- ✅ **Social login ready** for when you add social identity providers

## 📁 Directory Structure

```
infrastructure/sam/ui-customization/
├── README.md                      # Detailed documentation
├── QUICKSTART.md                  # This file
├── custom.css                     # Custom stylesheet
├── logo.svg                       # SARC-NG logo
├── upload-ui-assets.sh           # Upload to S3 (optional)
├── apply-ui-customization.sh     # Apply to Cognito
└── ui-assets-urls-*.txt          # Generated URLs (if using S3)
```

## 🚀 Quick Start (2 Steps)

### Step 1: Deploy Cognito Stack

First, make sure your Cognito stack is deployed:

```bash
cd /workspaces/sarc-ng/infrastructure/sam
make deploy-cognito ENV=dev
```

Wait for the deployment to complete (~2-3 minutes).

### Step 2: Apply UI Customization

Apply your custom login UI:

```bash
make ui-customize ENV=dev
```

That's it! Your custom login UI is now live! 🎉

## 🧪 Testing Your Custom UI

Get your Hosted UI login URL:

```bash
# Get the URL
aws cloudformation describe-stacks \
  --stack-name sarc-ng-cognito-dev \
  --query 'Stacks[0].Outputs[?OutputKey==`CognitoHostedUIUrl`].OutputValue' \
  --output text

# Or use the urls command
make urls ENV=dev
```

Open the URL in your browser. You should see:
- Your custom logo
- Purple gradient background
- Modern, clean login form
- Smooth animations

**Note:** Changes may take 2-5 minutes to propagate. Clear your browser cache (Ctrl+Shift+R) if you don't see changes immediately.

## 🎨 Customizing the Design

### Change Brand Colors

Edit `ui-customization/custom.css` and modify the CSS variables at the top:

```css
:root {
  --primary-color: #0066cc;      /* Your brand color */
  --primary-hover: #0052a3;      /* Darker shade for hover */
  /* ... other colors ... */
}
```

### Change Background Gradient

Find the `body` selector in `custom.css`:

```css
body {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  /* Change to your colors: */
  /* background: linear-gradient(135deg, #YOUR_COLOR_1 0%, #YOUR_COLOR_2 100%); */
}
```

### Replace Logo

Replace `logo.svg` with your own logo:
- **Recommended format:** SVG (scales perfectly)
- **Also supported:** PNG or JPG
- **Max size:** 2 MB
- **Recommended dimensions:** 200x80 pixels

### Apply Changes

After editing CSS or logo:

```bash
make ui-customize ENV=dev
```

Wait 2-5 minutes and refresh your browser.

## 📋 Available Make Commands

```bash
# Apply UI customization (recommended)
make ui-customize ENV=dev

# Upload assets to S3 (optional, for long-term URL storage)
make ui-upload ENV=dev

# Apply customization from local files
make ui-apply ENV=dev
```

## 🔍 What Gets Customized?

The custom CSS styles these Cognito elements:

| Element | Customization |
|---------|---------------|
| Login form | Modern card with shadow, rounded corners |
| Input fields | Clean borders, focus states, validation colors |
| Buttons | Brand colors, hover effects, loading states |
| Logo | Your custom logo at the top |
| Background | Gradient background (customizable) |
| Error messages | Styled alerts with icons |
| Links | Brand-colored with hover effects |
| Loading states | Custom spinners |
| MFA inputs | Styled code input boxes |

## 🌓 Dark Mode

The UI automatically supports dark mode! If a user's system is set to dark mode, they'll see:
- Dark background
- Light text
- Adjusted colors for readability

No configuration needed—it just works!

## 📱 Responsive Design

The login UI looks great on all devices:
- **Mobile:** Single column, touch-friendly buttons
- **Tablet:** Optimized spacing and sizing
- **Desktop:** Full experience with all features

## ⚙️ Advanced: Using S3 for Assets

If you want to host your CSS and logo on S3 (useful for referencing from external tools or long-term URLs):

```bash
# Upload to S3
make ui-upload ENV=dev

# This creates:
# - S3 bucket: sarc-ng-ui-assets-dev
# - CSS at: s3://bucket/cognito-ui/custom.css
# - Logo at: s3://bucket/cognito-ui/logo.svg
# - Pre-signed URLs saved to: ui-assets-urls-dev.txt
```

Then you can reference these URLs in AWS Console if needed.

## 🐛 Troubleshooting

### Changes Not Showing

1. **Wait 2-5 minutes** for CDN propagation
2. **Clear browser cache:** Ctrl+Shift+R (Windows/Linux) or Cmd+Shift+R (Mac)
3. **Try incognito/private mode**
4. **Check browser console** for errors

### Error: "Stack not found"

Deploy the Cognito stack first:
```bash
make deploy-cognito ENV=dev
```

### Error: "CSS file too large"

The CSS must be under 100 KB. Current size is ~20 KB, so this shouldn't happen unless you add a lot.

### Logo Not Showing

- Check file size (must be < 2 MB)
- Verify file format (SVG, PNG, or JPG)
- Make sure the file is named `logo.svg` (or update the script)

### IAM Permission Errors

You need these IAM permissions:
```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "cognito-idp:SetUICustomization",
      "cognito-idp:GetUICustomization",
      "cloudformation:DescribeStacks"
    ],
    "Resource": "*"
  }]
}
```

## 📚 More Information

- **Full Documentation:** See `README.md` in this directory
- **AWS Docs:** [Cognito UI Customization](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-ui-customization.html)
- **CSS Reference:** [Cognito CSS Classes](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-ui-customization.html#cognito-user-pools-app-ui-customization-css)

## 🎯 Next Steps

After setting up your custom login UI:

1. **Test the login flow:** Create a test user and sign in
2. **Customize further:** Adjust colors, fonts, and spacing
3. **Set up your app:** Integrate the Hosted UI with your application
4. **Enable MFA:** Add multi-factor authentication for security
5. **Add social login:** Configure Google, Facebook, or other providers

## 💡 Tips

- Keep a backup of your custom CSS
- Test on multiple devices and browsers
- Use the browser's developer tools to preview CSS changes
- Consider creating different styles for dev/staging/prod
- Check the login flow on mobile devices

## 🎨 Example Customizations

### Corporate Blue Theme
```css
--primary-color: #003d82;
--primary-hover: #002a5c;
body {
  background: linear-gradient(135deg, #005bbb 0%, #0078d4 100%);
}
```

### Green/Nature Theme
```css
--primary-color: #28a745;
--primary-hover: #218838;
body {
  background: linear-gradient(135deg, #56ab2f 0%, #a8e063 100%);
}
```

### Dark/Professional Theme
```css
--primary-color: #6c757d;
--primary-hover: #5a6268;
body {
  background: linear-gradient(135deg, #2c3e50 0%, #34495e 100%);
}
```

## 🎉 You're Done!

Your Cognito Hosted UI is now customized and ready to use. Enjoy your beautiful, branded login experience!

For questions or issues, see the full `README.md` or check AWS Cognito documentation.
