# build.js

**build.js** will be a generic tool to compile large JavaScript projects and manage dependencies both within projects and across projects.

# Overview

As I developed more and more JavaScript libraries for the web, I grew to use a very specific mechanism for managing dependencies and isolating code. All of the source files in my JavaScript libraries are not enclosed by functions. In addition, the code uses an `exports` variable to expose functions and values to other APIs or to other files within the API. At compile time, each source file is wrapped in a function and given an object for the `exports` variable. The exports variable is traditionally accessed globally as `window.APINAME`, creating a reasonable way for APIs to access each other on the web.

# LICENSE

**build.js** is licensed under the Don't Be A Dick Public License.

> DON'T BE A DICK PUBLIC LICENSE
> TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

 1. Do whatever you like with the original work, just don't be a dick.

     Being a dick includes - but is not limited to - the following instances:

	 1a. Outright copyright infringement - Don't just copy this and change the name.  
	 1b. Selling the unmodified original with no work done what-so-ever, that's REALLY being a dick.  
	 1c. Modifying the original work to contain hidden harmful content. That would make you a PROPER dick.  

 2. If you become rich through modifications, related works/services, or supporting the original work,
 share the love. Only a dick would make loads off this work and not buy the original work's 
 creator(s) a pint.
 
 3. Code is provided with no warranty. Using somebody else's code and bitching when it goes wrong makes 
 you a DONKEY dick. Fix the problem yourself. A non-dick would submit the fix back.
