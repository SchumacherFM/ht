Tutorial for Writing Test with ht
=================================

This tutorial describes everything one needs to know to use cmd/ht and write
tests and to combine tests into suites.  The tutorial is intended to be used
in interactive trainings; it may seem overly verbose if read silently.


Running ht
----------

`ht` is a command line application which does need external infastructure,
no Java/Ruby/Python runtime, no special libraries, no configuration files,
no registry entries, nothing. Download the version for your operating
system from (github)[https://github.com/vdobler/ht/releases].

Invoking it without any arguments describes it's usage and shows the
available subcommands. Help for the subcommands can be displayed in the
obvious way:

    $ ./ht
 
    $ ./ht help run

Make sure you use the proper version:
 
    $ ./ht version


Writing Tests
-------------

A test is stored as a JSON object in a file. The object needs three fields:

    {
        "Name":    "Some descriptive name, but no fancy characters please",
        "Request": { ... },
        "Checks":  [ ... ]
    }

Note the field names start with a Capital Letter and will be CamelCase.
Note that acronyms will be in all caps, e.g. "URL"

 * `Name` is a string an is needed to display and log the test properly.
 * `Request` is an object and contains information about the request to make
    for this test.
 * `Checks` is an array if checks to perform on the received response.


### JSON5 actually

The files are not JSON but (kinda [json5 comments]) JSON5: You are allowed
to omit the quotation marks `"` around fields names (if they are "normal"
field names) and the last element in an object or an array may have a trailing
comma `,` and you may add comments:


    {
        // This is a line comment
        Name:    "Some descriptive name, but no fancy characters please",
        Request: { ... },
        Checks:  [ ... ],
    }


### The Request URL

The main (and mandatory) field of `Request` is the `URL` which must be a
complete URL including schema, host (optional port) and path. The URL may
contain a fragment, but this won't be sent to the server.

    {
        Name:    "Unic homepage german",
        Request: {
            URL: "https://www.unic.com/de.html",
        },
        Checks:  [ ],
    }

TLS (SSL) secured request can be made, https:// is supported.
Save this as unic-hp1.ht and execute the request:

    $ ./ht run unic-hp1.ht

This should work, i.e. produce a request and print some output indicating
success but it is almost pointless as no checks are done.


### Checks

Checks provide high- (some low-) level checks on the received response.
We'll start with low-level checks as these are easier to understand.

    {
        Name:    "Unic homepage german",
        Request: {
            URL: "https://www.unic.com/de.html",
        },
        Checks:  [
            { Check: "StatusCode", Expect: 200 },
            { Check: "Body", Contains: "Unic" },
        ],
    }

Which check to execute is given in the `Check` field. Note again the
CamelCase nameing scheme.  A list of available checks is given below.
The other fields determine details of the check and are check dependent.
These other fields names have been chosen to allow "reading" the
check definition almost as clear text:

 * Check the status code and expect a value of 200.
 * Check the request body, it must contain "Unic".

You may want to run it (after saving to unic-hp2.ht):

    $ ./ht run unic-hp2.ht

You should see the passed checks.



Details of the `Request` object
-------------------------------

In the following we will take a detailed look at all fields of the `Request`
object and how to fine-control the generated request, how to send parameters,
how to add headers and cookies to the requests, etc. pp.


### POST, HEAD, PUT, DELETE...

The default for a request is `GET`.  If you want to create a different
type of request just specify the `Method`:

    {
        Name:    "Unic homepage german",
        Request: {
            Method: "HEAD",   //  Note: ALLCAPS as in the actual request
            URL: "https://www.unic.com/de.html",
        },
        Checks:  [
            { Check: "StatusCode", Expect: 200 },
            { Check: "Body", Contains: "Unic" },
        ],
    }

If you run this test you will see that the second checks fails: Well the
body of a HEAD request is empty and does not contain "Unic".

If `Method` is unset it defaults to `GET`.


### Sending query parameters

Sending parameters is quite simple as `ht` does all the heavy lifting of
encoding the parameters:

    {
        Name:    "Unic search",
        Request: {
            URL: "https://www.unic.com/de/tools/suche.html",
            Params: {
                q: "Magento",
                foo: "a b+c%d",
            }
        },
        Checks:  [
            { Check: "StatusCode", Expect: 200 },
            { Check: "Body", Contains: "Magento" },
        ],
    }

When running this test you will see that the value of the parameter `foo` is
properly encoded and that the parameters are sent as query parameters in
the URL.


### Sending POST parameters

How the parameters are sent is controlled with the `ParamsAs` field of `Request`:

    {
        Name:    "Unic search",
        Request: {
            Method: "POST",
            URL: "https://www.unic.com/de/tools/suche.html",
            ParamsAs: "body",  //  -->  application/x-www-form-urlencoded
            Params: {
                q: "Magento",
                foo: "a b+c%d",
            }
        },
        Checks:  [
            { Check: "StatusCode", Expect: 200 },
            { Check: "Body", Contains: "Magento" },
        ],
    }

This will make a POST-request and send the parameters urlencoded in the request
body. It will automatically set the appropriate Content-Type header to
"application/x-www-form-urlencoded".

Running this will fail as your Adobe AEM does not allow POST requests on
this URL.

If `ParamsAs` is unset it defaults to `URL` which indicates to send as query
parameters in the URL.


### Sending Parameters in body *and* URL

This is possible, all you have to to is encode the URL query parameter yourself:

    {
        Name:    "Unic search",
        Request: {
            Method: "POST",
            URL: "https://www.unic.com/de/tools/suche.html?a=b&c=%2B",  // c=+
            ParamsAs: "body",  
            Params: { q: "Magento" }
        },
        Checks:  [
            { Check: "StatusCode", Expect: 200 },
            { Check: "Body", Contains: "Magento" },
        ],
    }

This is not limited to POST request, this method works also for GET and other
requests.


### Sending files or general multipart forms

Ht can do multipart/form-data POST request like shown below. Such requests allow
file uploads for which `ht` has a special `@file`-syntax:

    {
        Name:    "Unic search",
        Request: {
            Method: "POST",
            URL: "https://www.unic.com/de/tools/suche.html",
            ParamsAs: "multipart",  
            Params: {
                q: "Magento",
                fileparam: "@file:path/to/file/to/upload",
                inlinefile: "@file:@name-to-send:the-data-to-send",
            }
        },
        Checks:  [
            { Check: "StatusCode", Expect: 200 },
            { Check: "Body", Contains: "Magento" },
        ],
    }

The inline form is useful if you want to POST e.g. a tiny JSON file which
can be typed into the test while the explicit form reads e.g. a larger image
file from the filesystem. 


### Sending multiple values for one parameter

This is possible, all you have to do is give the full set of values as a JSON
array:

    {
        Name:    "Unic search",
        Request: {
            URL: "https://www.unic.com/de/tools/suche.html",
            Params: {
                q: [ "Magento", "Sitecore", "hybris", "Drupal" ],
            }
        },
        Checks:  [ ... ]
    }


### Setting HTTP headers

The `Request` has a field `Header` which allows to set arbitrary header values:

    {
        Name:    "Unic homepage german",
        Request: {
            URL: "https://www.unic.com/de.html",
            Header: {
    	        User-Agent: [ "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.101 Safari/537.36" ],
    	        Accept: [ "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8" ],
    	        Accept-Language: [ "fr;q=0.2" ],
    	        Accept-Encoding: [ "gzip, deflate" ]
            }
        },
        Checks:  [
            { Check: "StatusCode", Expect: 200 },
            { Check: "Body", Contains: "Unic" },
        ],
    }

Note that the header fields are JSON arrays; no short syntax for setting a
header just once like it is available for parameters is available here.


### Following redirects

Sometimes you are interested only in the final response of a redirect chain:
You request some page, the server redirects you once or twice and on the
final request a page is delivered to you.  If you are not interested in the
redirect chain but solely on the outcome you may set `FollowRedirects`

    {
        Name:    "Unic homepage german",
        Request: {
            URL: "http://www.unic.com/",
            FollowRedirects: true,
        },
        Checks:  [
            { Check: "StatusCode", Expect: 200 },
            { Check: "Body", Contains: "Unic" },
        ],
    }

Running this will show you both requests made, one to `http://www.unic.com/` and
one to `https://www.unic.com/de.html`. The checks run on the response to the 
final  request.


### Sending cookies

You may send cookies through a handcrafted `Header` value (see above) or
with the convenience field `Cookies` in the `Request` object:

    {
        Name:    "Unic homepage german",
        Request: {
            URL: "http://www.unic.com/",
            Cookies: [
                {Name: "cip", Value: "602414252.20480.0000"},
                {Name: "WT_FPC", Value: "id=ad374f68-d3ba-4f9c-b351-1e58de56da56:lv=1442825498354:ss=1442825498354"},
            ]
        },
        Checks:  [ ... ]
    }


### Sending a body

If you want to send your own body (e.g. in a POST request) you set the `Body` field
of the `Request` object to the string you want to send.

Please consider this functionality as experimental.


Mixins (BasedOn)
----------------


It is painful and error prone to add the common header fields of a "normal"
browser-like request. To facilitate this `ht` provides the possibility to
merge partial tests -- call "mixins" -- into the actual test.

### Including partial tests

Assume you have the following JSON5 file which is a partial test (as it
has no URL):

    std-headers.mixin:
    {
        Name:    "Standard headers",
        Request: {
            Header: {
    	        User-Agent: [ "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.101 Safari/537.36" ],
    	        Accept: [ "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8" ],
    	        Accept-Language: [ "fr;q=0.2" ],
    	        Accept-Encoding: [ "gzip, deflate" ]
            }
        },
    }

and a "real test":

    real-test.ht
    {
        Name:    "Unic homepage german",
        BasedOn: [ "std-headers.mixin" ],  //  <<-- here
        Request: { URL: "https://www.unic.com/de.html" },
        Checks:  [ {Check:"StatusCode",Expect:200}, {Check:"Body",Contains: "Unic"},
        ],
    }

The settings made in std-headers.mixin are incorporated into the real-test.ht as
if a Header field would be present. The mixins are merged, the real-test.ht
himself may have Header definitions.


### Inclusion works transitively

Mixins them self can use mixins:

    std-webpage.ht
    {
        Name:    "Sensible settings and checks for each webpage",
        BasedOn: [ "../std-headers.mixin" ],  //  <<-- here
        Checks:  [
            { Check: "StatusCode", Expect: 200 },
            { Check: "UTF8Encoded" },
            { Check: "ContentType", Is: "text/html"},
            { Check: "ResponseTime", Lower: "2.5s" },
            { Check: "Body", Contains: "© {{NOW | "2006"}} Unic" },
        ],
    }

The mixin std-webpage.ht itself includes std-headers.mixin and add several
checks.  No any test may be `BasedOn` this std-webpage.ht and it will have
proper headers as sent by a browser and do four common checks on the response
automatically.

As you see relative filenames work too.



Combining Tests to Suites
-------------------------

Running several tests in one batch is possible by

    $ ./ht run test1.ht test2.ht test3.ht

But these tests are executed completely unrelated. For more control tests
can be combined into "Suites". 

Suites are stored on disc as JSON5 files like tests are. The following
shows everything a suite may contain.  As usual it starts with a Name
and as Description field:

    {
        Name:        "Sample Suite",
        Description: "Optional verbose details for suite",
        KeepCookies: true,   //  handle cookies like a browser
        OmitChecks:  false,  //  allows to do the requests but skip checking
        Verbosity:   2,      //  fix verbosity of all tests to 2

        Setup: [
          "reddit-homepage.ht"
        ],

        Tests: [
            "unic-search.ht",
            "google-homepage.ht",
            "unic-logo.ht",
            "unic-homepage.ht",
            "heise-homepage.ht",
            "sz-homepage.ht",
            "amazon-autocomplete.ht",
            "redirect.ht",
        ],

        Teardown: [
           "reddit-golang.ht",
           "reddit-programming.ht",
        ],

        Variables: {
           HOST: "www.unic.com",
           FOOBAR: "Something else here"
        }
    }   

The Setup, Tests and Teardown are arrays of filenames of tests.
The actual tests in Tests are executed only if all tests in Setup pass.
Teardown tests are executed always but their status is ignored.

All tests are executed strictly in serial order, one after the other.
If `KeepCookies` is true than any cookie set by the server will be
stored and (depending on the request details) sent back in subsequent
requests.

All this -- with the exception of the `Variables` field -- should be pretty
straightforward and obvious how it works.


Executing suites
----------------

Suites are nice but they offer structure which you might want to control.
The general way to execute a suite is to run the suite through `exec`:

    $ ./ht exec _somefancy.suite_

You may run several suites in one batch, `ht exec` will execute all suites
given in the command line.

To see which tests are "in" a suite use the `list` subcommand:

    $ ./ht list the.suite

If you want to skip certain tests or run just some test you can use the flags
`-skip` and `-only`. E.g.:

    $ ./ht exec -only 3-9 -skip 6 somefancy.suite

Would run only the actual (non-setup, non-teardown) tests 3, 4, 5, 7, 8 and 9.
(Counting from 1).


Using Variables in tests
------------------------

The most fancy part in the suite above is the `Variables` field which contains
key/value-pairs: Variable names and the corresponding value.  


### Using variables

Variable replacement can be used in a lot of places, from the request URL, over
parameter values to fields in checks. Variable replacements are written like

    "This is fixed {{VARNAME}} rest is fixed too."

If `VARNAME="foo 123" the resulting string will be:

    "This is fixed foo 123 rest is fixed too."

This works basically just like using ${VARNAME} in bash.
(Variables may have lowercase letters too.)

Take a look at the example suite above, 
`HOST` is a good example why variables exist:  You may want to write _one_ test
and have this test executed accessing different environments: From
localhost to development, to integration, to acceptance and even on
production.  Making the tests parametrized on the HOST name makes this
possible,


### Setting variables from the command line

Will I have to write suites for every environment just to provide variable
values? Of course not. `ht` has two command line flags which allow to set
variable values during invocation of `ht`:

 * `-D `_varname_`=`_value_ : Will set the variable _varname_ to _value_.
   E.g. `-D HOST=localhost:9001`

 * `-Dfile` _file.json_ : This will read variable names and values from the
   given JSON5 file _file.json_.

The `-Dfile` flags are handled first, you can overwrite the values with `-D`.

    $ ./ht -D HOST=127.0.0.1:8080 
 

### Replacing integer values with variable value

Yes, this can be done, even if not recommended. Just a small example:

    { Check: "Body", Contains: "Foo", Count: 9999 }

Checks that body contains exactly 9999 instances of "Foo". You can replace
this 9999 with a variable amount by defining a variable `#9999` (that is the
name!) and value 4.

This is a deadly hack, just think what you'll be testing when running a
suite with `-D "#200=900"`. But this hack its use when unrolling tests (see
table driven tests below).




Details of a Test
-----------------

A test has more fields than just the three (Name, Request, Checks) already
discussed. Some are useful when crafting a test.


### Retrying or "polling" and URL

Sometimes a failure is anticipated, e.g. while you wait for a server to start.
For such cases you may retry a tests several times.  This retrying is done
through the `Poll` field of a test: AN URL is polled until success:

    {
        Name:    "Unic homepage german",
        Request: { URL: "https://www.unic.com/de/html" },
        Poll: {
            Max: 12,
            Sleep: "7654ms",
        },
        Checks:  [ {Check: "StatusCode", Expect: 200} ],
    }

This test will be done up to 12 times. If it passes earlier the test succeeds
and if it fails for all 12 tries it fails.


### Debug related stuff

You can keep more information about the test in the Description fields (e.g.
a reference to a requirement or a bug number) and you can set the verbosity
during tests execution:

    {
        Name:        "Unic homepage german",
        Description: "More information goes here".
        Request:     { URL: "https://www.unic.com/de.html" },
        Checks:      [{Check:"StatusCode",Expect:200},{Check:"Body",Contains:"Unic"}],
        Verbosity:   3,       //  3 --> TRACE
    }



### Timeouts and Sleep

Tests have some more fields which control timing:

    {
        Name:        "Unic homepage german",
        Timeout:     "5.3s",
        Request:     { URL: "https://www.unic.com/de.html" },
        Checks:      [{Check:"StatusCode",Expect:200},{Check:"Body",Contains:"Unic"}],
        PreSleep:    0.5,     //  number means seconds,
        InterSleep:  "9ms"    //  strings are durations
        PostSleep:   "3m15s"  //     "     "      " 
    }

If you want tot set a different timeout for the requests (dial, send request,
read full response) you use the `Timeout` field.  Note that the Timeout field
is part of Test, not Request.

The `Pre`-, `Inter`- and `Post`-Sleep fields control how much time is slept
before starting the test, between the request and the checks and after the
checks.  There usefulness is not obvious at this stage of the tutorial.



Table driven Tests
------------------

Test can be used as a template and this templated can be unrolled resulting in
several individual tests, each a bit differently parametrized.
Currently this works only if a test is run as part of a suite (so not with
the `run` subcommand).


### Simple unrolling

This unrolling of the given test-template is done by the `Unroll` object
and uses variable substitution. E.g.

    {
        Name:    "Unic homepage german",
        Request: {
            Method: "{{METHOD}}",
            URL: "https://www.unic.com/de.html",
        },
        Checks:  [
            { Check: "StatusCode", Expect: 200 },
            { Check: "Body", Contains: "Unic" },
        ],
        Unroll: {
            METHOD: [ "GET", "HEAD", "POST" ],
        }
    }

If such a test(-template) is included in a suite definition, the resulting
suite will have 3 test at the place of this template: One with a GET method,
one with HEAD and one with POST.  The last two will fail.

Unrolling works on more than one variable too and this is how the
test can be made successful:

    {
        Name:    "Unic homepage german",
        Request: {
            Method: "{{METHOD}}",
            URL: "https://www.unic.com/de.html",
        },
        Checks:  [
            { Check: "StatusCode", Expect: 999 }, // <-- 999 !
            { Check: "Body", Contains: "{{TEXT}}" },
        ],
        Unroll: {
            METHOD: [ "GET",  "HEAD", "POST" ],
            TEXT:   [ "Unic", "",     "Forbidden" ],
            "#999": [ "200",  "200",  "403" ],
        }
    }

Note that the text to look for is empty fro the HEAD request and "Forbidden"
for the POST request.  For the status code which is an integer (and you cannot
use normal variables with integers) the ugly hack from above is used.

This section was labeled "simple" unrolling because all three variables are
unrolled through three different values each.  If some variable has a different
number of values it is no longer simple.


### Fancy unrolling

Let's take a look at an synthectic example:

    {
        Name:    "Syntetic example for fancy unrolling",
        Request: {
            URL: "http://www.example.org/{{AAA}}/{{BBB}}.html",
        },
        Unroll: {
            AAA: [ "java", "php" ],
            BBB: [ "hybris", "magento", "sitecore" ],
        },
    }

This would be unrolled to 6 = 2 values for AAA * 3 values for BBB tests with
the following paths beeing requested:

    /java/hybris.html
    /java/magento.html
    /java/sitecore.html
    /php/hybris.html
    /php/magento.html
    /php/sitecore.html

As 2 and 3 are relatively prime (well they are actually prime :-) each of the
possible 6 combinations is used.

If the number of unrolled variable values is not relatively prime, the number
of generated test will be the least common multiple. Examples might make this
clear:

One is a factor of the other like here:

    Unroll {
        AAA: [ "java", "php" ],
        BBB: [ "hybris", "magento", "sitecore", "aem6" ],
    }

would produce just 4 combinations: hybris/java, magento/php, sitecore/java and
aem6/php.

The following has 4 (=2*2) and 6 (=2*3) values resulting 2*2*3=12 test:
    Unroll {
        AAA: [ "a", "b", "c", "d" ], 
        BBB: [ "1", "2", "3", "4", "5", "6" ],
    }

a/1, b/2, c/3, d/4, a/5, b/6, c/1, d/2, a/3, b/4, c/5 and d/6.


Normally you will just want one of the following to:

 *  All unrolled variables have the same number of values:
    Values unroll in lockstep.

 *  You use relatively prime numbers and get the full set of combinations.  

















[json5 comments]: The files are very JSON5 like, but the JSON5 parser
used has one bug: Comments may not occur after the last element of
a slice (or an object). 