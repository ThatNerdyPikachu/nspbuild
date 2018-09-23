# nspbuild
**nspbuild** is a tool that automates the process of creating homebrew NSPs for the Nintendo Switch.
For the uninformed, NSPs are pretty much the CIAs (installable titles) of the Switch.

This tool aims to help anyone who wants to create their own NSP, without having information scattered around forums and such.

# The Guide™

**Note: To follow this guide, you'll need a populated keys.txt in the same directory as nspbuild**

## The NSO
Let's say we want to create an NSP of [Checkpoint](https://github.com/FlagBrew/Checkpoint) by [Bernardo Giordano](https://github.com/BernardoGiordano). First, we'd ``git clone`` it, then compile it. We'll get a few files, but the one we're interested in is the ``.nso`` file. We can keep that in a safe place, feel free to delete everything else.

## The Icon
(You can skip this step by specifying ``none`` as the icon, but who wants that? Nevertheless, you can do this if you'd like.)

Time to get out your MS Paint skills, let's get designing! For this, you'll need a 256x256 JPEG with **zero exif data**

To make sure your image has zero exif data, open it with your favorite image editor, save it as a BMP, open that BMP, then resave it as a JPEG.

## Actually creating our NSP
Sweet, we're almost done! Here's the final step, actually building our NSP! What we want to do now is download the [latest release](https://github.com/ThatNerdyPikachu/nspbuild/releases/latest) of nspbuild, and put it in a folder with our NSO and our icon. For this example, let's say our icon is named ``icon.jpg``, and our NSO is named ``Checkpoint.nso``. When we run the program with zero arguments, we get a help message:
```
usage: nspbuild <path/to/nso> <name> <author> <version> <path/to/icon/jpg> <tid>
```

So, that's exactly what we'll specify!
```
nspbuild Checkpoint.nso Checkpoint "Bernardo Giordano" 3.4.2 icon.jpg 01005791048912af
```

Once it's done, we should see a new folder called ``out/``! In it, what's this? ``Checkpoint [01005791048912af].nsp``!

# Troubleshooting
If you're having any issues with the program, join my [Discord server](https://invite.gg/pika) and ask for help! We'll be happy to assist you!

# Credits
None of this wouldn't have been possible without these amazing people!
```
The-4n, for creating hacBuildPack
roblabla, for creating linkle
switchbrew, for creating switch-tools, and subsequently npdmtool
The Golang Authors, for creating such an amazing language!
```