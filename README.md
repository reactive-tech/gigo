
[Gigo](https://gigo.reactive-tech.io/) is a static website generator written in Go allowing web developers to write 
websites in core HTML, JS and CSS while benefiting from features missing in HTML 5 such as including and reusing pages. 
When using Gigo there is no need to learn a custom templating language. 
Just work with the standards core HTML, JS and CSS and write portable websites.

HTML specification does not offer a way of including pages and Gigo resolves this simply and elegantly.

Gigo has the following features:

* Tag &lt;gigo-include file=* /&gt; : includes other HTML pages in the page using this tag.

* Tag &lt;gigo-include-in file=* /&gt; : the page using this tag asks to be included in another template page, reducing 
  further code duplications.
  
* It is written in Go and comes with a binary which is self-contained and does not require the installation of a third 
  party library or SDK.

* Because it does not rely on a database, you can commit your work in GIT and generate anytime a static website using 
  gigo binary.

[Gigo](https://gigo.reactive-tech.io/) was developed by [Reactive Tech Limited](https://www.reactive-tech.io/)  and Alex Arica as the lead developer.

Get-started:
https://gigo.reactive-tech.io/doc/getting-started.html
